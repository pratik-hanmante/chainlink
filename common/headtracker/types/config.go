package types

import "time"

type Config interface {
	BlockEmissionIdleWarningThreshold() time.Duration
	FinalityDepth() uint32
	FinalityTag() bool
}

type HeadTrackerConfig interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
	TotalHeadsLimit() int
	PersistHeads() bool
	CanonicalChainRule() string
}
