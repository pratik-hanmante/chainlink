package headtracker

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"

	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type headSaver struct {
	heads         map[common.Hash]*evmtypes.Head
	mu            sync.RWMutex
	canonicalHead *evmtypes.Head
	finalizedHead *evmtypes.Head
	orm           ORM
	htConfig      HeadTrackerConfig
	logger        logger.Logger
}

func NewHeadSaver(lggr logger.Logger, orm ORM, htConfig HeadTrackerConfig) httypes.HeadSaver {
	return &headSaver{}
}

func (hs *headSaver) SaveHead(ctx context.Context, head *evmtypes.Head, canonical bool, latestFinalized *evmtypes.Head) (err error) {
	if hs.htConfig.PersistHeads() {
		if err = hs.orm.IdempotentInsertHead(ctx, head); err != nil {
			return err
		}
		if err = hs.orm.TrimOldHeads(ctx, uint(latestFinalized.Number)); err != nil {
			return err
		}
	}

	hs.saveHeads(latestFinalized, canonical, head)

	return nil
}

func (hs *headSaver) LatestCanonicalHead() *evmtypes.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.canonicalHead
}

func (hs *headSaver) LatestFinalizedHead() *evmtypes.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.finalizedHead
}

func (hs *headSaver) Chain(hash common.Hash) *evmtypes.Head {
	hs.mu.RLock()
	defer hs.mu.RUnlock()

	return hs.heads[hash]
}

func (h *headSaver) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.heads)
}

func (hs *headSaver) LoadHeads(ctx context.Context, latestFinalized *evmtypes.Head) (int, error) {
	heads, err := hs.orm.LatestHeads(ctx, uint(latestFinalized.BlockNumber()))
	if err != nil {
		return 0, err
	}

	hs.saveHeads(latestFinalized, false, heads...)
	return hs.Count(), nil
}

func (hs *headSaver) LatestHeadFromDB(ctx context.Context) (head *evmtypes.Head, err error) {
	return hs.orm.LatestHead(ctx)
}


func (hs *headSaver) saveHeads(finalized *evmtypes.Head, canonicalHead bool, newHeads ...*evmtypes.Head) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	headsMap := make(map[common.Hash]*evmtypes.Head, len(hs.heads)+len(newHeads)+1)
	headsMap[finalized.Hash] = finalized
	finalized.Finalized = true
	for _, newHead := range newHeads {
		if newHead.Number > finalized.Number {
			headsMap[newHead.Hash] = newHead
		}
	}

	for _, head := range hs.heads {
		if head.Hash == head.ParentHash {
			// shouldn't happen but it is untrusted input
			continue
		}
		// copy all head objects to avoid races when a previous head chain is used
		// elsewhere (since we mutate Parent here)
		headCopy := *head
		headCopy.Parent = nil // always build it from scratch in case it points to a head too old to be included
		headCopy.Finalized = false
		// Map eliminates duplicates. Include every unfinalized block up until the latest finalized.
		if headCopy.Number > finalized.Number {
			headsMap[head.Hash] = &headCopy
		}
	}

	for _, head := range hs.heads {
		parent, exists := headsMap[head.ParentHash]
		if exists {
			head.Parent = parent
		}
	}

	if canonicalHead && len(newHeads) == 1 && headsMap[newHeads[0].Hash] != nil {
		hs.canonicalHead = headsMap[newHeads[0].Hash]
	}
	hs.finalizedHead = headsMap[finalized.Hash]

	// set
	hs.heads = headsMap

	if hs.Count() > hs.htConfig.TotalHeadsLimit() {
		hs.logger.Warnw("chain larger than EvmHeadTrackerTotalHeadsLimit. In memory heads exceed limit.",
			"headsCount", hs.Count(), "evmHeadTrackerTotalHeadsLimit", hs.htConfig.TotalHeadsLimit())
	}
}

var NullSaver httypes.HeadSaver = &nullSaver{}

type nullSaver struct{}

func (*nullSaver) SaveHead(ctx context.Context, head *evmtypes.Head, canonical bool, latestFinalized *evmtypes.Head) error {
	return nil
}
func (*nullSaver) LoadHeads(ctx context.Context, latestFinalized *evmtypes.Head) (int, error) {
	return 0, nil
}
func (*nullSaver) LatestHeadFromDB(ctx context.Context) (*evmtypes.Head, error) { return nil, nil }
func (*nullSaver) LatestCanonicalHead() *evmtypes.Head                          { return nil }
func (*nullSaver) LatestFinalizedHead() *evmtypes.Head                          { return nil }
func (*nullSaver) Chain(hash common.Hash) *evmtypes.Head                        { return nil }
