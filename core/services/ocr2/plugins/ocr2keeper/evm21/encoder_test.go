package evm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/stretchr/testify/assert"
)

func TestEVMAutomationEncoder21(t *testing.T) {
	encoder := EVMAutomationEncoder21{}

	t.Run("encoding an empty list of upkeep results returns a nil byte array", func(t *testing.T) {
		b, err := encoder.Encode()
		assert.Equal(t, ErrEmptyResults, err)
		assert.Equal(t, b, []byte(nil))
	})

	// t.Run("successfully encodes a single upkeep result", func(t *testing.T) {
	// 	upkeepResult := EVMAutomationUpkeepResult21{
	// 		Block:            1,
	// 		ID:               big.NewInt(10),
	// 		Eligible:         true,
	// 		GasUsed:          big.NewInt(100),
	// 		PerformData:      []byte("data"),
	// 		FastGasWei:       big.NewInt(100),
	// 		LinkNative:       big.NewInt(100),
	// 		CheckBlockNumber: 1,
	// 		CheckBlockHash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	// 		ExecuteGas:       10,
	// 	}
	// 	b, err := encoder.EncodeReport([]ocr2keepers.UpkeepResult{upkeepResult})
	// 	assert.Nil(t, err)
	// 	assert.Len(t, b, 416)

	// 	t.Run("successfully decodes a report with a single upkeep result", func(t *testing.T) {
	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Nil(t, err)
	// 		assert.Len(t, upkeeps, 1)

	// 		upkeep := upkeeps[0].(EVMAutomationUpkeepResult21)

	// 		// some fields aren't populated by the decode so we compare field-by-field for those that are populated
	// 		assert.Equal(t, upkeep.Block, upkeepResult.Block)
	// 		assert.Equal(t, upkeep.ID, upkeepResult.ID)
	// 		assert.Equal(t, upkeep.Eligible, upkeepResult.Eligible)
	// 		assert.Equal(t, upkeep.PerformData, upkeepResult.PerformData)
	// 		assert.Equal(t, upkeep.FastGasWei, upkeepResult.FastGasWei)
	// 		assert.Equal(t, upkeep.LinkNative, upkeepResult.LinkNative)
	// 		assert.Equal(t, upkeep.CheckBlockNumber, upkeepResult.CheckBlockNumber)
	// 		assert.Equal(t, upkeep.CheckBlockHash, upkeepResult.CheckBlockHash)
	// 	})

	// 	t.Run("an error is returned when unpacking into a map fails", func(t *testing.T) {
	// 		oldUnpackIntoMapFn := unpackIntoMapFn
	// 		unpackIntoMapFn = func(v map[string]interface{}, data []byte) error {
	// 			return errors.New("failed to unpack into map")
	// 		}
	// 		defer func() {
	// 			unpackIntoMapFn = oldUnpackIntoMapFn
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "failed to unpack into map")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when an expected key is missing from the map", func(t *testing.T) {
	// 		oldMKeys := mKeys
	// 		mKeys = []string{"fastGasWei", "linkNative", "upkeepIds", "wrappedPerformDatas", "thisKeyWontExist"}
	// 		defer func() {
	// 			mKeys = oldMKeys
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "decoding error")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when the third element of the map is not a slice of big.Int", func(t *testing.T) {
	// 		oldMKeys := mKeys
	// 		mKeys = []string{"fastGasWei", "linkNative", "wrappedPerformDatas", "upkeepIds"}
	// 		defer func() {
	// 			mKeys = oldMKeys
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "upkeep ids of incorrect type in report")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when the fourth element of the map is not a struct of perform data", func(t *testing.T) {
	// 		oldMKeys := mKeys
	// 		mKeys = []string{"fastGasWei", "linkNative", "upkeepIds", "upkeepIds"}
	// 		defer func() {
	// 			mKeys = oldMKeys
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "performs of incorrect structure in report")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when the upkeep ids and performDatas are of different lengths", func(t *testing.T) {
	// 		oldUnpackIntoMapFn := unpackIntoMapFn
	// 		unpackIntoMapFn = func(v map[string]interface{}, data []byte) error {
	// 			v["fastGasWei"] = 1
	// 			v["linkNative"] = 2
	// 			v["upkeepIds"] = []*big.Int{big.NewInt(123), big.NewInt(456)}
	// 			v["wrappedPerformDatas"] = []struct {
	// 				CheckBlockNumber uint32   `json:"checkBlockNumber"`
	// 				CheckBlockhash   [32]byte `json:"checkBlockhash"`
	// 				PerformData      []byte   `json:"performData"`
	// 			}{
	// 				{
	// 					CheckBlockNumber: 1,
	// 					CheckBlockhash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	// 					PerformData:      []byte{},
	// 				},
	// 			}
	// 			return nil
	// 		}
	// 		defer func() {
	// 			unpackIntoMapFn = oldUnpackIntoMapFn
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "upkeep ids and performs should have matching length")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when the first element of the map is not a big int", func(t *testing.T) {
	// 		oldMKeys := mKeys
	// 		mKeys = []string{"upkeepIds", "linkNative", "upkeepIds", "wrappedPerformDatas"}
	// 		defer func() {
	// 			mKeys = oldMKeys
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "fast gas as wrong type")
	// 		assert.Len(t, upkeeps, 0)
	// 	})

	// 	t.Run("an error is returned when the second element of the map is not a big int", func(t *testing.T) {
	// 		oldMKeys := mKeys
	// 		mKeys = []string{"fastGasWei", "upkeepIds", "upkeepIds", "wrappedPerformDatas"}
	// 		defer func() {
	// 			mKeys = oldMKeys
	// 		}()

	// 		upkeeps, err := encoder.DecodeReport(b)
	// 		assert.Error(t, err, "link native as wrong type")
	// 		assert.Len(t, upkeeps, 0)
	// 	})
	// })

	// t.Run("successfully encodes multiple upkeep results", func(t *testing.T) {
	// 	upkeepResult0 := EVMAutomationUpkeepResult21{
	// 		Block:            1,
	// 		ID:               big.NewInt(10),
	// 		Eligible:         true,
	// 		GasUsed:          big.NewInt(100),
	// 		PerformData:      []byte("data0"),
	// 		FastGasWei:       big.NewInt(100),
	// 		LinkNative:       big.NewInt(100),
	// 		CheckBlockNumber: 1,
	// 		CheckBlockHash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	// 		ExecuteGas:       10,
	// 	}
	// 	upkeepResult1 := EVMAutomationUpkeepResult21{
	// 		Block:            1,
	// 		ID:               big.NewInt(10),
	// 		Eligible:         true,
	// 		GasUsed:          big.NewInt(200),
	// 		PerformData:      []byte("data1"),
	// 		FastGasWei:       big.NewInt(200),
	// 		LinkNative:       big.NewInt(200),
	// 		CheckBlockNumber: 2,
	// 		CheckBlockHash:   [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	// 		ExecuteGas:       20,
	// 	}
	// 	b, err := encoder.EncodeReport([]ocr2keepers.UpkeepResult{upkeepResult0, upkeepResult1})
	// 	assert.Nil(t, err)
	// 	assert.Len(t, b, 640)
	// })

	t.Run("an error is returned when pack fails", func(t *testing.T) {
		oldPackFn := packFn
		packFn = func(args ...interface{}) ([]byte, error) {
			return nil, errors.New("pack failed")
		}
		defer func() {
			packFn = oldPackFn
		}()

		result := ocr2keepers.CheckResult{
			Payload: ocr2keepers.UpkeepPayload{
				Upkeep: ocr2keepers.ConfiguredUpkeep{
					ID: ocr2keepers.UpkeepIdentifier([]byte("10")),
				},
				Trigger: ocr2keepers.Trigger{
					BlockNumber: 1,
					BlockHash:   common.Bytes2Hex([]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}),
				},
			},
			Eligible:     true,
			GasAllocated: 100,
			PerformData:  []byte("data0"),
			Extension: EVMAutomationResultExtension21{
				FastGasWei: big.NewInt(100),
				LinkNative: big.NewInt(100),
			},
		}

		b, err := encoder.Encode(result)
		assert.Errorf(t, err, "pack failed: failed to pack report data")
		assert.Len(t, b, 0)
	})

}
