// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_allow_list

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
)

var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

var TermsOfServiceAllowListMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCallableByRouterOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RecipientIsBlocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RouterMustBeSet\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"enabled\",\"type\":\"bool\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"acceptTermsOfService\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"blockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAllowedSenders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getConfigHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"config\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"messageHash\",\"type\":\"bytes32\"}],\"name\":\"getEthSignedMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"acceptor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"getMessageHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"hasAccess\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"isBlockedSender\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"unblockSender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"config\",\"type\":\"bytes\"}],\"name\":\"updateConfig\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620012ce380380620012ce833981016040819052620000349162000164565b81816001600160a01b0382166200005e57604051632530e88560e11b815260040160405180910390fd5b600180546001600160a01b0319166001600160a01b038416179055620000848162000099565b80516020909101206000555062000292915050565b60008082806020019051810190620000b291906200024f565b6040805180820182528315158082526001600160a01b0384166020928301819052600580546001600160a81b031916610100600160a81b031984161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b6001600160a01b03811681146200014b57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b600080604083850312156200017857600080fd5b8251620001858162000135565b602084810151919350906001600160401b0380821115620001a557600080fd5b818601915086601f830112620001ba57600080fd5b815181811115620001cf57620001cf6200014e565b604051601f8201601f19908116603f01168101908382118183101715620001fa57620001fa6200014e565b8160405282815289868487010111156200021357600080fd5b600093505b8284101562000237578484018601518185018701529285019262000218565b60008684830101528096505050505050509250929050565b600080604083850312156200026357600080fd5b825180151581146200027457600080fd5b6020840151909250620002878162000135565b809150509250929050565b61102c80620002a26000396000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806382184c7b116100815780639883c10d1161005b5780639883c10d146101b1578063a5e1d61d146101b9578063fa540801146101cc57600080fd5b806382184c7b146101785780638cc6acce1461018b57806395a078e81461019e57600080fd5b80636b14daf8116100b25780636b14daf81461011f57806380e8a15114610142578063817ef62e1461016357600080fd5b8063181f5a77146100d95780632179d447146100f757806347663acb1461010c575b600080fd5b6100e161022d565b6040516100ee9190610bb3565b60405180910390f35b61010a610105366004610c8d565b61024d565b005b61010a61011a366004610cf2565b6103f2565b61013261012d366004610d0f565b610535565b60405190151581526020016100ee565b610155610150366004610d64565b61055f565b6040519081526020016100ee565b61016b6105bd565b6040516100ee9190610d9d565b61010a610186366004610cf2565b6105ce565b61010a610199366004610e26565b610720565b6101326101ac366004610cf2565b610787565b600054610155565b6101326101c7366004610cf2565b6107a7565b6101556101da366004610ef5565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c8101829052600090605c01604051602081830303815290604052805190602001209050919050565b6060604051806060016040528060288152602001610ff860289139905090565b73ffffffffffffffffffffffffffffffffffffffff831660009081526004602052604090205460ff16156102ad576040517f62b7a34d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600554610100900473ffffffffffffffffffffffffffffffffffffffff166103176102db6101da878761055f565b84848080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506107e892505050565b73ffffffffffffffffffffffffffffffffffffffff1614610364576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff84161415806103a957503373ffffffffffffffffffffffffffffffffffffffff8516148015906103a95750333b155b156103e0576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103eb6002846108bf565b5050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af1158015610461573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104859190610f0e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104e9576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff16600090815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b60055460009060ff1661054a57506001610558565b6105556002856108e1565b90505b9392505050565b6040517fffffffffffffffffffffffffffffffffffffffff000000000000000000000000606084811b8216602084015283901b1660348201526000906048016040516020818303038152906040528051906020012090505b92915050565b60606105c96002610910565b905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638da5cb5b6040518163ffffffff1660e01b81526004016020604051808303816000875af115801561063d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106619190610f0e565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106c5576040517fa0f0a44600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106d060028261091d565b5073ffffffffffffffffffffffffffffffffffffffff16600090815260046020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610771576040517fc41a5b0900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61077a8161093f565b8051602090910120600055565b60055460009060ff1661079c57506001919050565b6105b76002836108e1565b60055460009060ff166107bc57506000919050565b5073ffffffffffffffffffffffffffffffffffffffff1660009081526004602052604090205460ff1690565b6000806000808451604114610829576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050602082810151604080850151606080870151835160008082529681018086528a9052951a928501839052840183905260808401819052919260019060a0016020604051602081039080840390855afa15801561088c573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe00151979650505050505050565b60006105588373ffffffffffffffffffffffffffffffffffffffff8416610a15565b73ffffffffffffffffffffffffffffffffffffffff811660009081526001830160205260408120541515610558565b6060600061055883610a64565b60006105588373ffffffffffffffffffffffffffffffffffffffff8416610ac0565b600080828060200190518101906109569190610f2b565b60408051808201825283151580825273ffffffffffffffffffffffffffffffffffffffff84166020928301819052600580547fffffffffffffffffffffff000000000000000000000000000000000000000000167fffffffffffffffffffffff0000000000000000000000000000000000000000ff84161761010090920291909117905591519182529294509092507f22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f910160405180910390a1505050565b6000818152600183016020526040812054610a5c575081546001818101845560008481526020808220909301849055845484825282860190935260409020919091556105b7565b5060006105b7565b606081600001805480602002602001604051908101604052809291908181526020018280548015610ab457602002820191906000526020600020905b815481526020019060010190808311610aa0575b50505050509050919050565b60008181526001830160205260408120548015610ba9576000610ae4600183610f5f565b8554909150600090610af890600190610f5f565b9050818114610b5d576000866000018281548110610b1857610b18610f99565b9060005260206000200154905080876000018481548110610b3b57610b3b610f99565b6000918252602080832090910192909255918252600188019052604090208390555b8554869080610b6e57610b6e610fc8565b6001900381819060005260206000200160009055905585600101600086815260200190815260200160002060009055600193505050506105b7565b60009150506105b7565b600060208083528351808285015260005b81811015610be057858101830151858201604001528201610bc4565b5060006040828601015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b73ffffffffffffffffffffffffffffffffffffffff81168114610c4157600080fd5b50565b60008083601f840112610c5657600080fd5b50813567ffffffffffffffff811115610c6e57600080fd5b602083019150836020828501011115610c8657600080fd5b9250929050565b60008060008060608587031215610ca357600080fd5b8435610cae81610c1f565b93506020850135610cbe81610c1f565b9250604085013567ffffffffffffffff811115610cda57600080fd5b610ce687828801610c44565b95989497509550505050565b600060208284031215610d0457600080fd5b813561055881610c1f565b600080600060408486031215610d2457600080fd5b8335610d2f81610c1f565b9250602084013567ffffffffffffffff811115610d4b57600080fd5b610d5786828701610c44565b9497909650939450505050565b60008060408385031215610d7757600080fd5b8235610d8281610c1f565b91506020830135610d9281610c1f565b809150509250929050565b6020808252825182820181905260009190848201906040850190845b81811015610deb57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101610db9565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060208284031215610e3857600080fd5b813567ffffffffffffffff80821115610e5057600080fd5b818401915084601f830112610e6457600080fd5b813581811115610e7657610e76610df7565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610ebc57610ebc610df7565b81604052828152876020848701011115610ed557600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208284031215610f0757600080fd5b5035919050565b600060208284031215610f2057600080fd5b815161055881610c1f565b60008060408385031215610f3e57600080fd5b82518015158114610f4e57600080fd5b6020840151909250610d9281610c1f565b818103818111156105b7577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fdfe46756e6374696f6e73205465726d73206f66205365727669636520416c6c6f77204c697374207631a164736f6c6343000813000a",
}

var TermsOfServiceAllowListABI = TermsOfServiceAllowListMetaData.ABI

var TermsOfServiceAllowListBin = TermsOfServiceAllowListMetaData.Bin

func DeployTermsOfServiceAllowList(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address, config []byte) (common.Address, *types.Transaction, *TermsOfServiceAllowList, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TermsOfServiceAllowListBin), backend, router, config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TermsOfServiceAllowList{TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

type TermsOfServiceAllowList struct {
	address common.Address
	abi     abi.ABI
	TermsOfServiceAllowListCaller
	TermsOfServiceAllowListTransactor
	TermsOfServiceAllowListFilterer
}

type TermsOfServiceAllowListCaller struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListTransactor struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListFilterer struct {
	contract *bind.BoundContract
}

type TermsOfServiceAllowListSession struct {
	Contract     *TermsOfServiceAllowList
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListCallerSession struct {
	Contract *TermsOfServiceAllowListCaller
	CallOpts bind.CallOpts
}

type TermsOfServiceAllowListTransactorSession struct {
	Contract     *TermsOfServiceAllowListTransactor
	TransactOpts bind.TransactOpts
}

type TermsOfServiceAllowListRaw struct {
	Contract *TermsOfServiceAllowList
}

type TermsOfServiceAllowListCallerRaw struct {
	Contract *TermsOfServiceAllowListCaller
}

type TermsOfServiceAllowListTransactorRaw struct {
	Contract *TermsOfServiceAllowListTransactor
}

func NewTermsOfServiceAllowList(address common.Address, backend bind.ContractBackend) (*TermsOfServiceAllowList, error) {
	abi, err := abi.JSON(strings.NewReader(TermsOfServiceAllowListABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindTermsOfServiceAllowList(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowList{address: address, abi: abi, TermsOfServiceAllowListCaller: TermsOfServiceAllowListCaller{contract: contract}, TermsOfServiceAllowListTransactor: TermsOfServiceAllowListTransactor{contract: contract}, TermsOfServiceAllowListFilterer: TermsOfServiceAllowListFilterer{contract: contract}}, nil
}

func NewTermsOfServiceAllowListCaller(address common.Address, caller bind.ContractCaller) (*TermsOfServiceAllowListCaller, error) {
	contract, err := bindTermsOfServiceAllowList(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListCaller{contract: contract}, nil
}

func NewTermsOfServiceAllowListTransactor(address common.Address, transactor bind.ContractTransactor) (*TermsOfServiceAllowListTransactor, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListTransactor{contract: contract}, nil
}

func NewTermsOfServiceAllowListFilterer(address common.Address, filterer bind.ContractFilterer) (*TermsOfServiceAllowListFilterer, error) {
	contract, err := bindTermsOfServiceAllowList(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListFilterer{contract: contract}, nil
}

func bindTermsOfServiceAllowList(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TermsOfServiceAllowListMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListCaller.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.TermsOfServiceAllowListTransactor.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TermsOfServiceAllowList.Contract.contract.Call(opts, result, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transfer(opts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.contract.Transact(opts, method, params...)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getAllAllowedSenders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetAllAllowedSenders() ([]common.Address, error) {
	return _TermsOfServiceAllowList.Contract.GetAllAllowedSenders(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetConfigHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getConfigHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetConfigHash() ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetConfigHash(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetConfigHash() ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetConfigHash(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetEthSignedMessageHash(opts *bind.CallOpts, messageHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getEthSignedMessageHash", messageHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetEthSignedMessageHash(messageHash [32]byte) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetEthSignedMessageHash(&_TermsOfServiceAllowList.CallOpts, messageHash)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetEthSignedMessageHash(messageHash [32]byte) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetEthSignedMessageHash(&_TermsOfServiceAllowList.CallOpts, messageHash)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) GetMessageHash(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "getMessageHash", acceptor, recipient)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) GetMessageHash(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessageHash(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) GetMessageHash(acceptor common.Address, recipient common.Address) ([32]byte, error) {
	return _TermsOfServiceAllowList.Contract.GetMessageHash(&_TermsOfServiceAllowList.CallOpts, acceptor, recipient)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "hasAccess", user, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) HasAccess(user common.Address, arg1 []byte) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess(&_TermsOfServiceAllowList.CallOpts, user, arg1)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) HasAccess0(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "hasAccess0", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) HasAccess0(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess0(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) HasAccess0(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.HasAccess0(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "isBlockedSender", sender)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) IsBlockedSender(sender common.Address) (bool, error) {
	return _TermsOfServiceAllowList.Contract.IsBlockedSender(&_TermsOfServiceAllowList.CallOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCaller) TypeAndVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TermsOfServiceAllowList.contract.Call(opts, &out, "typeAndVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListCallerSession) TypeAndVersion() (string, error) {
	return _TermsOfServiceAllowList.Contract.TypeAndVersion(&_TermsOfServiceAllowList.CallOpts)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "acceptTermsOfService", acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) AcceptTermsOfService(acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.AcceptTermsOfService(&_TermsOfServiceAllowList.TransactOpts, acceptor, recipient, proof)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "blockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) BlockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.BlockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "unblockSender", sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UnblockSender(sender common.Address) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UnblockSender(&_TermsOfServiceAllowList.TransactOpts, sender)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactor) UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.contract.Transact(opts, "updateConfig", config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListTransactorSession) UpdateConfig(config []byte) (*types.Transaction, error) {
	return _TermsOfServiceAllowList.Contract.UpdateConfig(&_TermsOfServiceAllowList.TransactOpts, config)
}

type TermsOfServiceAllowListConfigSetIterator struct {
	Event *TermsOfServiceAllowListConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *TermsOfServiceAllowListConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TermsOfServiceAllowListConfigSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}

	select {
	case log := <-it.logs:
		it.Event = new(TermsOfServiceAllowListConfigSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

func (it *TermsOfServiceAllowListConfigSetIterator) Error() error {
	return it.fail
}

func (it *TermsOfServiceAllowListConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type TermsOfServiceAllowListConfigSet struct {
	Enabled bool
	Raw     types.Log
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &TermsOfServiceAllowListConfigSetIterator{contract: _TermsOfServiceAllowList.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error) {

	logs, sub, err := _TermsOfServiceAllowList.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(TermsOfServiceAllowListConfigSet)
				if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowListFilterer) ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error) {
	event := new(TermsOfServiceAllowListConfigSet)
	if err := _TermsOfServiceAllowList.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _TermsOfServiceAllowList.abi.Events["ConfigSet"].ID:
		return _TermsOfServiceAllowList.ParseConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (TermsOfServiceAllowListConfigSet) Topic() common.Hash {
	return common.HexToHash("0x22aa8545955b447cb49ea37e67de742e750839c633ded8c9b5b09614843b229f")
}

func (_TermsOfServiceAllowList *TermsOfServiceAllowList) Address() common.Address {
	return _TermsOfServiceAllowList.address
}

type TermsOfServiceAllowListInterface interface {
	GetAllAllowedSenders(opts *bind.CallOpts) ([]common.Address, error)

	GetConfigHash(opts *bind.CallOpts) ([32]byte, error)

	GetEthSignedMessageHash(opts *bind.CallOpts, messageHash [32]byte) ([32]byte, error)

	GetMessageHash(opts *bind.CallOpts, acceptor common.Address, recipient common.Address) ([32]byte, error)

	HasAccess(opts *bind.CallOpts, user common.Address, arg1 []byte) (bool, error)

	HasAccess0(opts *bind.CallOpts, sender common.Address) (bool, error)

	IsBlockedSender(opts *bind.CallOpts, sender common.Address) (bool, error)

	TypeAndVersion(opts *bind.CallOpts) (string, error)

	AcceptTermsOfService(opts *bind.TransactOpts, acceptor common.Address, recipient common.Address, proof []byte) (*types.Transaction, error)

	BlockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UnblockSender(opts *bind.TransactOpts, sender common.Address) (*types.Transaction, error)

	UpdateConfig(opts *bind.TransactOpts, config []byte) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*TermsOfServiceAllowListConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *TermsOfServiceAllowListConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*TermsOfServiceAllowListConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
