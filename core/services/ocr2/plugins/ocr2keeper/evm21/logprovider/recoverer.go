package logprovider

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	keepersflows "github.com/smartcontractkit/ocr2keepers/pkg/v3/flows"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type UpkeepStateReader interface {
	// SelectByUpkeepIDsAndBlockRange retrieves upkeep states for provided upkeep ids and block range, the result is currently not in particular order
	SelectByUpkeepIDsAndBlockRange(upkeepIds []*big.Int, start, end int64) ([]*ocr2keepers.UpkeepPayload, []*ocr2keepers.UpkeepState, error)
}

type RecoveryOptions struct {
	Interval       time.Duration
	GCInterval     time.Duration
	TTL            time.Duration
	LookbackBlocks int64
	BatchSize      int32
}

func (o *RecoveryOptions) defaults() {
	if o.Interval == 0 {
		o.Interval = 30 * time.Second
	}
	if o.GCInterval == 0 {
		o.GCInterval = o.Interval*10 + o.Interval/2
	}
	if o.TTL == 0 {
		o.TTL = 24*time.Hour - time.Second
	}
	if o.LookbackBlocks == 0 {
		o.LookbackBlocks = 512
	}
	if o.BatchSize == 0 {
		o.BatchSize = 10
	}
}

type logRecoverer struct {
	lggr logger.Logger

	cancel context.CancelFunc

	lock *sync.RWMutex

	opts RecoveryOptions

	pending []ocr2keepers.UpkeepPayload

	filterStore *upkeepFilterStore
	visited     map[string]time.Time

	poller       logpoller.LogPoller
	upkeepStates UpkeepStateReader
}

var _ keepersflows.RecoverableProvider = &logRecoverer{}

func NewLogRecoverer(lggr logger.Logger, poller logpoller.LogPoller, upkeepStates UpkeepStateReader, opts *RecoveryOptions) *logRecoverer {
	if opts == nil {
		opts = new(RecoveryOptions)
	}
	opts.defaults()
	filterStore := newUpkeepFilterStore()
	return &logRecoverer{
		lggr:        lggr.Named("LogRecoverer"),
		opts:        *opts,
		lock:        &sync.RWMutex{},
		pending:     make([]ocr2keepers.UpkeepPayload, 0),
		filterStore: filterStore,
		visited:     make(map[string]time.Time),
		poller:      poller,
	}
}

func (r *logRecoverer) Start(pctx context.Context) error {
	ctx, cancel := context.WithCancel(pctx)
	r.lock.Lock()
	r.cancel = cancel
	interval := r.opts.Interval
	gcInterval := r.opts.GCInterval
	r.lock.Unlock()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	gcTicker := time.NewTicker(gcInterval)
	defer gcTicker.Stop()

	r.lggr.Debug("Starting log recoverer")

	for {
		select {
		case <-ticker.C:
			r.recover(ctx)
		case <-gcTicker.C:
			r.clean(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *logRecoverer) Close() error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *logRecoverer) GetRecoverables() ([]ocr2keepers.UpkeepPayload, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.pending) == 0 {
		return nil, nil
	}

	pending := r.pending
	r.pending = make([]ocr2keepers.UpkeepPayload, 0)

	for _, p := range pending {
		r.visited[p.ID] = time.Now()
	}

	return pending, nil
}

func (r *logRecoverer) clean(ctx context.Context) {
	r.lock.Lock()
	defer r.lock.Unlock()

	cleaned := 0
	for id, t := range r.visited {
		if time.Since(t) > r.opts.TTL {
			delete(r.visited, id)
			cleaned++
		}
	}

	if cleaned > 0 {
		r.lggr.Debugw("gc: cleaned visited upkeeps", "cleaned", cleaned)
	}
}

func (r *logRecoverer) recover(ctx context.Context) error {
	latest, err := r.poller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrHeadNotAvailable, err)
	}
	offsetBlock := r.getRecoveryOffsetBlock(latest)
	if offsetBlock < 0 {
		// too soon to recover, we don't have enough blocks
		return nil
	}

	filters := r.getFiltersBatch(offsetBlock)
	if len(filters) == 0 {
		return nil
	}

	r.lggr.Debugw("recovering logs", "filters", filters)

	var wg sync.WaitGroup
	for _, f := range filters {
		wg.Add(1)
		go func(f upkeepFilter) {
			defer wg.Done()
			r.recoverFilter(ctx, f)
		}(f)
	}
	wg.Wait()

	return nil
}

func (r *logRecoverer) recoverFilter(ctx context.Context, f upkeepFilter) error {
	start, end := f.lastRecoveredBlock, f.lastRecoveredBlock+10
	logs, err := r.poller.LogsWithSigs(start, end, f.eventSigs, f.contract, pg.WithParentCtx(ctx))
	if err != nil {
		return fmt.Errorf("could not read logs: %w", err)
	}
	payloads, _, err := r.upkeepStates.SelectByUpkeepIDsAndBlockRange([]*big.Int{f.upkeepID}, start, end)
	if err != nil {
		return fmt.Errorf("could not read upkeep states: %w", err)
	}
	for _, log := range logs {
		trigger := logToTrigger(f.upkeepID, log)
		for _, payload := range payloads {
			if !compareTriggers(trigger, payload.Trigger) {
				continue
			}
			// checks if we already visited this log
			if ts, ok := r.visited[payload.ID]; ok && time.Since(ts) < r.opts.TTL {
				continue
			}
			r.lggr.Debugw("recovered missed log", "trigger", trigger, "upkeepID", f.upkeepID)
			r.pending = append(r.pending, *payload)
		}
	}
	return nil
}

func (r *logRecoverer) getFiltersBatch(offsetBlock int64) []upkeepFilter {
	filters := make([]upkeepFilter, 0)

	r.filterStore.Range(func(f upkeepFilter) bool {
		if f.lastRecoveredBlock >= offsetBlock {
			return true
		}
		for i, filter := range filters {
			if filter.lastRecoveredBlock > f.lastRecoveredBlock {
				if i == 0 {
					filters = append([]upkeepFilter{f}, filters...)
				} else {
					filters = append(filters[:i], append([]upkeepFilter{f}, filters[i:]...)...)
				}
				return true
			}
		}
		filters = append(filters, f)
		return true
	})

	return r.selectBatch(filters)
}

func (r *logRecoverer) selectBatch(filters []upkeepFilter) []upkeepFilter {
	batchSize := int(atomic.LoadInt32(&r.opts.BatchSize))

	if len(filters) < batchSize {
		return filters
	}
	results := filters[:batchSize/2]
	filters = filters[batchSize/2:]

	for len(results) < batchSize && len(filters) != 0 {
		i := rand.Intn(len(filters))
		results = append(results, filters[i])
		if i == 0 {
			filters = filters[1:]
		} else if i == len(filters)-1 {
			filters = filters[:i]
		} else {
			filters = append(filters[:i], filters[i+1:]...)
		}
	}

	return results
}

func (r *logRecoverer) getRecoveryOffsetBlock(latest int64) int64 {
	lookbackBlocks := atomic.LoadInt64(&r.opts.LookbackBlocks)
	return latest - lookbackBlocks
}

func logToTrigger(id *big.Int, log logpoller.Log) ocr2keepers.Trigger {
	return ocr2keepers.NewTrigger(
		log.BlockNumber,
		log.BlockHash.Hex(),
		LogTriggerExtension{
			TxHash:   log.TxHash.Hex(),
			LogIndex: log.LogIndex,
		},
	)
}

func compareTriggers(t1, t2 ocr2keepers.Trigger) bool {
	return t1.BlockNumber == t2.BlockNumber &&
		t1.BlockHash == t2.BlockHash &&
		t1.Extension == t2.Extension // TODO: compare extension
}
