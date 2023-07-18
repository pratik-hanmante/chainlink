package headtracker

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/maps"

	htrktypes "github.com/smartcontractkit/chainlink/v2/common/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

var (
	promCurrentHead = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "head_tracker_current_head",
		Help: "The highest seen head number",
	}, []string{"evmChainID"})

	promOldHead = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "head_tracker_very_old_head",
		Help: "Counter is incremented every time we get a head that is much lower than the highest seen head ('much lower' is defined as a block that is EVM.FinalityDepth or greater below the highest seen head)",
	}, []string{"evmChainID"})
)

// HeadsBufferSize - The buffer is used when heads sampling is disabled, to ensure the callback is run for every head
const HeadsBufferSize = 10

type HeadTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
] struct {
	log             logger.Logger
	headBroadcaster types.HeadBroadcaster[HTH, BLOCK_HASH]
	headSaver       types.HeadSaver[HTH, BLOCK_HASH]
	mailMon         *utils.MailboxMonitor
	client          htrktypes.Client[HTH, S, ID, BLOCK_HASH]
	chainID         ID
	config          htrktypes.Config
	htConfig        htrktypes.HeadTrackerConfig

	backfillMB   *utils.Mailbox[HTH]
	broadcastMB  *utils.Mailbox[HTH]
	headListener types.HeadListener[HTH, BLOCK_HASH]
	chStop       utils.StopChan
	wgDone       sync.WaitGroup
	utils.StartStopOnce
	getNilHead func() HTH
}

// NewHeadTracker instantiates a new HeadTracker using HeadSaver to persist new block numbers.
func NewHeadTracker[
	HTH htrktypes.Head[BLOCK_HASH, ID],
	S types.Subscription,
	ID types.ID,
	BLOCK_HASH types.Hashable,
](
	lggr logger.Logger,
	client htrktypes.Client[HTH, S, ID, BLOCK_HASH],
	config htrktypes.Config,
	htConfig htrktypes.HeadTrackerConfig,
	headBroadcaster types.HeadBroadcaster[HTH, BLOCK_HASH],
	headSaver types.HeadSaver[HTH, BLOCK_HASH],
	mailMon *utils.MailboxMonitor,
	getNilHead func() HTH,
) types.HeadTracker[HTH, BLOCK_HASH] {
	chStop := make(chan struct{})
	lggr = lggr.Named("HeadTracker")
	return &HeadTracker[HTH, S, ID, BLOCK_HASH]{
		headBroadcaster: headBroadcaster,
		client:          client,
		chainID:         client.ConfiguredChainID(),
		config:          config,
		htConfig:        htConfig,
		log:             lggr,
		backfillMB:      utils.NewSingleMailbox[HTH](),
		broadcastMB:     utils.NewMailbox[HTH](HeadsBufferSize),
		chStop:          chStop,
		headListener:    NewHeadListener[HTH, S, ID, BLOCK_HASH](lggr, client, config, chStop),
		headSaver:       headSaver,
		mailMon:         mailMon,
		getNilHead:      getNilHead,
	}
}

// Start starts HeadTracker service.
func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) Start(ctx context.Context) error {
	return ht.StartOnce("HeadTracker", func() error {
		ht.log.Debugw("Starting HeadTracker", "chainID", ht.chainID)

		initialHead, err := ht.getInitialHead(ctx)
		if err != nil {
			if errors.Is(err, ctx.Err()) {
				return nil
			}
			return errors.Wrapf(err, "Error getting initial head")
		}

		if ht.htConfig.PersistHeads() {
			latestFinalized, err := ht.calculateLatestFinalized(ctx, initialHead.BlockNumber())
			if err != nil {
				return err
			}

			n, err := ht.headSaver.LoadHeads(ctx, latestFinalized)
			if err != nil {
				return err
			}
			ht.log.Debugw(
				fmt.Sprintf("HeadTracker: Loaded successfully %v blocks. Finalized block %s and block number %v",
					config.FriendlyNumber(n), latestFinalized.BlockHash(), latestFinalized.BlockNumber()))
		}

		if err := ht.handleNewHead(ctx, initialHead); err != nil {
			return errors.Wrap(err, "error handling initial head")
		}
		ht.wgDone.Add(2)
		go ht.headListener.ListenForNewHeads(ht.handleNewHead, ht.wgDone.Done)
		go ht.broadcastLoop()

		ht.mailMon.Monitor(ht.broadcastMB, "HeadTracker", "Broadcast", ht.chainID.String())

		return nil
	})
}

// Close stops HeadTracker service.
func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) Close() error {
	return ht.StopOnce("HeadTracker", func() error {
		close(ht.chStop)
		ht.wgDone.Wait()
		return ht.broadcastMB.Close()
	})
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) Name() string {
	return ht.log.Name()
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) HealthReport() map[string]error {
	report := map[string]error{
		ht.Name(): ht.StartStopOnce.Healthy(),
	}
	maps.Copy(report, ht.headListener.HealthReport())
	return report
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) LatestCanonicalChain() HTH {
	return ht.headSaver.LatestCanonicalHead()
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) Backfill(ctx context.Context, initialHead HTH, latestFinalized HTH) (err error) {
	head := initialHead.EarliestHeadInChain()
	headBlockNumber := head.BlockNumber()
	latestFinalizedBlockNumber := latestFinalized.BlockNumber()
	if headBlockNumber <= latestFinalized.BlockNumber() {
		return nil
	}
	mark := time.Now()
	fetched := 0
	l := ht.log.With("blockNumber", headBlockNumber,
		"n", headBlockNumber-latestFinalizedBlockNumber,
		"fromBlockHeight", latestFinalizedBlockNumber,
		"toBlockHeight", headBlockNumber-1)
	l.Debug("Starting backfill")
	defer func() {
		if ctx.Err() != nil {
			l.Warnw("Backfill context error", "err", ctx.Err())
			return
		}
		l.Debugw("Finished backfill",
			"fetched", fetched,
			"time", time.Since(mark),
			"err", err)
	}()

	for i := head.BlockNumber() - 1; i > latestFinalizedBlockNumber; i-- {
		// NOTE: Sequential requests here mean it's a potential performance bottleneck, be aware!
		existingHead := ht.headSaver.Chain(head.GetParentHash())
		if existingHead.IsValid() {
			head = existingHead
			continue
		}
		head, err = ht.fetchAndSaveHead(ctx, head.GetParentHash(), latestFinalized)
		if ctx.Err() != nil {
			ht.log.Debugw("context canceled, aborting backfill", "err", err, "ctx.Err", ctx.Err())
			break
		} else if err != nil {
			return errors.Wrap(err, "fetchAndSaveHead failed")
		}
		fetched++
	}
	if head.GetParentHash() != latestFinalized.GetParentHash() {
		return errors.New("Backfill failed: a reorg happened during backfill")
	}
	return
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) getInitialHead(ctx context.Context) (HTH, error) {
	head, err := ht.client.HeadByNumber(ctx, nil)
	if err != nil {
		return ht.getNilHead(), errors.Wrap(err, "failed to fetch initial head")
	}
	if head.IsValid() {
		loggerFields := []interface{}{"head", head}
		loggerFields = append(loggerFields, "blockNumber", head.BlockNumber(), "blockHash", head.BlockHash())
		ht.log.Debugw("Got initial head", loggerFields...)
		return head, nil
	} else {
		return ht.getNilHead(), errors.Wrap(err, "Got nil initial head")
	}
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) decideCanonicalChain(ctx context.Context, head HTH, latestFinalized HTH) (bool, error) {
	currentCanonical := ht.headSaver.LatestCanonicalHead()
	if !currentCanonical.IsValid() {
		ht.headSaver.SaveHead(ctx, head, true, latestFinalized)
		return true, nil
	}
	switch ht.htConfig.CanonicalChainRule() {
	case "PoS":
		err := ht.headSaver.SaveHead(ctx, head, true, latestFinalized)
		if ctx.Err() != nil {
			return false, nil
		} else if err != nil {
			return false, errors.Wrapf(err, "failed to save head: %#v", head)
		}
		return true, nil
	default:
		fallthrough
	case "LongestChain":
		storeCanonical := head.BlockNumber() > ht.headSaver.LatestCanonicalHead().BlockNumber()
		err := ht.headSaver.SaveHead(ctx, head, storeCanonical, latestFinalized)
		if ctx.Err() != nil {
			return false, nil
		} else if err != nil {
			return false, errors.Wrapf(err, "failed to save head: %#v", head)
		}

		return storeCanonical, nil
	}
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) handleNewHead(ctx context.Context, head HTH) error {
	latestFinalizedHead, err := ht.calculateLatestFinalized(ctx, head.BlockNumber())
	if err != nil {
		return err
	}

	if head.IsValid() {
		ht.log.Debugw(fmt.Sprintf("Received new head %v", config.FriendlyNumber(head.BlockNumber())),
			"blockHeight", head.BlockNumber(),
			"blockHash", head.BlockHash(),
			"parentHeadHash", head.GetParentHash(),
		)

	}

	if head.BlockNumber() <= latestFinalizedHead.BlockNumber() {
		if head.BlockNumber() == latestFinalizedHead.BlockNumber() && head.BlockHash() != latestFinalizedHead.BlockHash() {
			return errors.Wrapf(errors.New("Invariant violation: current finalized block has the same number as received head, but different hash!"),
				"blockNumber", head.BlockNumber(), "newHash", head.BlockHash(), "finalizedHash", latestFinalizedHead.BlockHash())
		}
		promOldHead.WithLabelValues(ht.chainID.String()).Inc()
		ht.log.Criticalf(`Got very old block with number %d (latest finalized %d).This is a problem and either means a very deep re-org occurred,`+
			`one of the RPC nodes has gotten far out of sync`, head.BlockNumber(), latestFinalizedHead.BlockNumber())
		ht.SvcErrBuffer.Append(errors.New("got very old block"))
		return nil
	}

	if ht.headSaver.Chain(head.BlockHash()).IsValid() {
		return nil
	}

	broadcast, err := ht.decideCanonicalChain(ctx, head, latestFinalizedHead)
	if err != nil {
		return err
	}
	if broadcast {
		promCurrentHead.WithLabelValues(ht.chainID.String()).Set(float64(head.BlockNumber()))
		headWithChain := ht.headSaver.Chain(head.BlockHash())
		if !headWithChain.IsValid() {
			return errors.Errorf("HeadTracker#handleNewHighestHead headWithChain was unexpectedly nil")
		}

		err = ht.Backfill(ctx, head, latestFinalizedHead)
		if err != nil {
			ht.log.Warnw("Unexpected error while backfilling heads", "err", err)
		}

		ht.broadcastMB.Deliver(headWithChain)
	}
	return nil
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) broadcastLoop() {
	defer ht.wgDone.Done()

	samplingInterval := ht.htConfig.SamplingInterval()
	if samplingInterval > 0 {
		ht.log.Debugf("Head sampling is enabled - sampling interval is set to: %v", samplingInterval)
		debounceHead := time.NewTicker(samplingInterval)
		defer debounceHead.Stop()
		for {
			select {
			case <-ht.chStop:
				return
			case <-debounceHead.C:
				item := ht.broadcastMB.RetrieveLatestAndClear()
				if !item.IsValid() {
					continue
				}
				ht.headBroadcaster.BroadcastNewLongestChain(item)
			}
		}
	} else {
		ht.log.Info("Head sampling is disabled - callback will be called on every head")
		for {
			select {
			case <-ht.chStop:
				return
			case <-ht.broadcastMB.Notify():
				for {
					item, exists := ht.broadcastMB.Retrieve()
					if !exists {
						break
					}
					ht.headBroadcaster.BroadcastNewLongestChain(item)
				}
			}
		}
	}
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) fetchAndSaveHead(ctx context.Context, headHash BLOCK_HASH, latestFinalized HTH) (HTH, error) {
	ht.log.Debugw("Fetching head", "blockHash", headHash)
	head, err := ht.client.HeadByHash(ctx, headHash)
	if err != nil {
		return ht.getNilHead(), err
	} else if !head.IsValid() {
		return ht.getNilHead(), errors.New("got nil head")
	}
	err = ht.headSaver.SaveHead(ctx, head, false, latestFinalized)
	if err != nil {
		return ht.getNilHead(), err
	}
	return head, nil
}

func (ht *HeadTracker[HTH, S, ID, BLOCK_HASH]) calculateLatestFinalized(ctx context.Context, currentHeadNumber int64) (HTH, error) {
	if !ht.config.FinalityTag() {
		// If heads is not empty, calculate the latest finalized head from the highest head saved, otherwise just use the current one
		var newFinalized int64
		previousFinalized := ht.headSaver.LatestFinalizedHead()
		if !previousFinalized.IsValid() {
			newFinalized = currentHeadNumber - int64(ht.config.FinalityDepth()-1)
		} else {
			newFinalized = mathutil.Max(previousFinalized.BlockNumber(), currentHeadNumber-int64(ht.config.FinalityDepth())-1)
			if newFinalized == previousFinalized.BlockNumber() {
				return previousFinalized, nil
			}
		}
		return ht.client.HeadByNumber(ctx, big.NewInt(mathutil.Max(newFinalized, 0)))
	}
	return ht.client.LatestBlockByType(ctx, "finalized")
}
