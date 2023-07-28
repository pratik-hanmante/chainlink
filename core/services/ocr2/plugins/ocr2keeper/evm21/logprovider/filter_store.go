package logprovider

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/time/rate"
)

type upkeepFilter struct {
	contract  common.Address
	eventSigs []common.Hash
	upkeepID  *big.Int
	// lastPollBlock is the last block number the logs were fetched for this upkeep
	// used by log event provider.
	lastPollBlock int64
	// blockLimiter is used to limit the number of blocks to fetch logs for an upkeep.
	// used by log event provider.
	blockLimiter *rate.Limiter
	// lastRePollBlock is the last block number the logs were recovered for this upkeep
	// used by log recoverer.
	lastRecoveredBlock int64
}

type upkeepFilterStore struct {
	lock    *sync.RWMutex
	filters map[string]upkeepFilter
}

func newUpkeepFilterStore() *upkeepFilterStore {
	return &upkeepFilterStore{
		lock:    &sync.RWMutex{},
		filters: map[string]upkeepFilter{},
	}
}

func (s *upkeepFilterStore) InitializeActiveUpkeeps(filters ...upkeepFilter) {
	s.lock.Lock()
	s.filters = make(map[string]upkeepFilter)
	s.lock.Unlock()

	s.AddActiveUpkeeps(filters...)
}

func (s *upkeepFilterStore) Range(iterator func(upkeepFilter) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for _, f := range s.filters {
		if !iterator(f) {
			break
		}
	}
}

func (s *upkeepFilterStore) GetIDs(iterator func(*big.Int) bool) []*big.Int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if iterator == nil {
		// noop iterator returns true for all filters
		iterator = func(*big.Int) bool { return true }
	}

	var ids []*big.Int
	for _, f := range s.filters {
		if iterator(f.upkeepID) {
			ids = append(ids, f.upkeepID)
		}
	}

	return ids
}

func (s *upkeepFilterStore) GetFilters(selector func(upkeepFilter) bool, ids ...*big.Int) []upkeepFilter {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if selector == nil {
		// noop selector returns true for all filters
		selector = func(upkeepFilter) bool { return true }
	}

	var filters []upkeepFilter
	for _, id := range ids {
		if f, ok := s.filters[id.String()]; ok && selector(f) {
			filters = append(filters, f)
		}
	}
	return filters
}

func (s *upkeepFilterStore) AddActiveUpkeeps(filters ...upkeepFilter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, f := range filters {
		s.filters[f.upkeepID.String()] = f
	}
}

func (s *upkeepFilterStore) RemoveActiveUpkeeps(filters ...upkeepFilter) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, f := range filters {
		uid := f.upkeepID.String()
		if _, ok := s.filters[uid]; ok {
			delete(s.filters, uid)
		}
	}
}
