// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package functions_client_example

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

var FunctionsClientExampleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"router\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EmptyArgs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptySource\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoInlineSecrets\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRouterCanFufill\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRequestID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"RequestSent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_CALLBACK_GAS\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"err\",\"type\":\"bytes\"}],\"name\":\"handleOracleFulfillment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastErrorLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastRequestId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponse\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_lastResponseLength\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"source\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"encryptedSecretsReferences\",\"type\":\"bytes\"},{\"internalType\":\"string[]\",\"name\":\"args\",\"type\":\"string[]\"},{\"internalType\":\"uint64\",\"name\":\"subscriptionId\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"jobId\",\"type\":\"bytes32\"}],\"name\":\"sendRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051620019e3380380620019e3833981016040819052620000349162000198565b600080546001600160a01b0319166001600160a01b038316178155339081906001600160a01b038216620000af5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600180546001600160a01b0319166001600160a01b0384811691909117909155811615620000e257620000e281620000ec565b50505050620001ca565b336001600160a01b03821603620001465760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401620000a6565b600280546001600160a01b0319166001600160a01b03838116918217909255600154604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b600060208284031215620001ab57600080fd5b81516001600160a01b0381168114620001c357600080fd5b9392505050565b61180980620001da6000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c80636d9809a011610081578063b1e217491161005b578063b1e2174914610182578063f2fde38b1461018b578063f7b4c06f1461019e57600080fd5b80636d9809a01461014857806379ba5097146101525780638da5cb5b1461015a57600080fd5b806342748b2a116100b257806342748b2a146100ff5780634b0795a81461012c5780635fa353e71461013557600080fd5b80630ca76175146100ce5780633944ea3a146100e3575b600080fd5b6100e16100dc366004611255565b6101ae565b005b6100ec60045481565b6040519081526020015b60405180910390f35b60065461011790640100000000900463ffffffff1681565b60405163ffffffff90911681526020016100f6565b6100ec60055481565b6100e1610143366004611328565b61020f565b6101176201117081565b6100e1610319565b60015460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f6565b6100ec60035481565b6100e161019936600461140c565b61041f565b6006546101179063ffffffff1681565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101ff576040517f5099014100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61020a838383610433565b505050565b610217610501565b6102586040805160e0810190915280600081526020016000815260200160008152602001606081526020016060815260200160608152602001606081525090565b61029a89898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105849050565b85156102e2576102e287878080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525085939250506105959050565b83156102fc576102fc6102f58587611442565b82906105df565b61030b81846201117085610622565b600355505050505050505050565b60025473ffffffffffffffffffffffffffffffffffffffff16331461039f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064015b60405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000008082163390811790935560028054909116905560405173ffffffffffffffffffffffffffffffffffffffff909116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a350565b610427610501565b61043081610646565b50565b8260035414610471576040517fd068bf5b00000000000000000000000000000000000000000000000000000000815260048101849052602401610396565b61047a8261073c565b6004558151600680547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9092169190911790556104bc8161073c565b600555516006805463ffffffff909216640100000000027fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff9092169190911790555050565b60015473ffffffffffffffffffffffffffffffffffffffff163314610582576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610396565b565b61059182600080846107be565b5050565b80516000036105d0576040517fe889636f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60016020830152608090910152565b805160000361061a576040517ffe936cb700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60a090910152565b60008061062e86610855565b905061063c81868686610bc2565b9695505050505050565b3373ffffffffffffffffffffffffffffffffffffffff8216036106c5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610396565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff838116918217909255600154604051919216907fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127890600090a350565b60008060209050602083511015610751575081515b60005b818110156107b757610767816008611538565b8482815181106107795761077961154f565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c92909217916107b08161157e565b9050610754565b5050919050565b80516000036107f9576040517f22ce3edd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8383600281111561080c5761080c6114da565b9081600281111561081f5761081f6114da565b90525060408401828015610835576108356114da565b90818015610845576108456114da565b9052506060909301929092525050565b606061085f61110c565b805161086d90610100610c9b565b506108ad816040518060400160405280600c81526020017f636f64654c6f636174696f6e0000000000000000000000000000000000000000815250610d15565b6108cc81846000015160028111156108c7576108c76114da565b610d2e565b61090b816040518060400160405280600881526020017f6c616e6775616765000000000000000000000000000000000000000000000000815250610d15565b61092581846040015160008111156108c7576108c76114da565b610964816040518060400160405280600681526020017f736f757263650000000000000000000000000000000000000000000000000000815250610d15565b610972818460600151610d15565b60a08301515115610a16576109bc816040518060400160405280600481526020017f6172677300000000000000000000000000000000000000000000000000000000815250610d15565b6109c581610d67565b60005b8360a0015151811015610a0c576109fc828560a0015183815181106109ef576109ef61154f565b6020026020010151610d15565b610a058161157e565b90506109c8565b50610a1681610d8b565b60808301515115610b1657600083602001516002811115610a3957610a396114da565b03610a70576040517fa80d31f700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610aaf816040518060400160405280600f81526020017f736563726574734c6f636174696f6e0000000000000000000000000000000000815250610d15565b610ac981846020015160028111156108c7576108c76114da565b610b08816040518060400160405280600781526020017f7365637265747300000000000000000000000000000000000000000000000000815250610d15565b610b16818460800151610da9565b60c08301515115610bba57610b60816040518060400160405280600981526020017f6279746573417267730000000000000000000000000000000000000000000000815250610d15565b610b6981610d67565b60005b8360c0015151811015610bb057610ba0828560c001518381518110610b9357610b9361154f565b6020026020010151610da9565b610ba98161157e565b9050610b6c565b50610bba81610d8b565b515192915050565b600080546040517f461d276200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff9091169063461d276290610c229087908990600190899089906004016115b6565b6020604051808303816000875af1158015610c41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c659190611656565b60405190915081907f1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db890600090a2949350505050565b604080518082019091526060815260006020820152610cbb60208361166f565b15610ce357610ccb60208361166f565b610cd69060206116aa565b610ce090836116bd565b91505b602080840183905260405180855260008152908184010181811015610d0757600080fd5b604052508290505b92915050565b610d228260038351610db2565b815161020a9082610ed9565b8151610d3b9060c2610f01565b506105918282604051602001610d5391815260200190565b604051602081830303815290604052610da9565b610d72816004610f6a565b600181602001818151610d8591906116bd565b90525050565b610d96816007610f6a565b600181602001818151610d8591906116aa565b610d2282600283515b60178167ffffffffffffffff1611610ddf578251610dd99060e0600585901b168317610f01565b50505050565b60ff8167ffffffffffffffff1611610e21578251610e08906018611fe0600586901b1617610f01565b508251610dd99067ffffffffffffffff83166001610f81565b61ffff8167ffffffffffffffff1611610e64578251610e4b906019611fe0600586901b1617610f01565b508251610dd99067ffffffffffffffff83166002610f81565b63ffffffff8167ffffffffffffffff1611610ea9578251610e9090601a611fe0600586901b1617610f01565b508251610dd99067ffffffffffffffff83166004610f81565b8251610ec090601b611fe0600586901b1617610f01565b508251610dd99067ffffffffffffffff83166008610f81565b604080518082019091526060815260006020820152610efa83838451611006565b9392505050565b6040805180820190915260608152600060208201528251516000610f268260016116bd565b905084602001518210610f4757610f4785610f42836002611538565b6110f5565b8451602083820101858153508051821115610f60578181525b5093949350505050565b815161020a90601f611fe0600585901b1617610f01565b6040805180820190915260608152600060208201528351516000610fa582856116bd565b90508560200151811115610fc257610fc286610f42836002611538565b60006001610fd2866101006117f0565b610fdc91906116aa565b90508651828101878319825116178152508051831115610ffa578281525b50959695505050505050565b604080518082019091526060815260006020820152825182111561102957600080fd5b835151600061103884836116bd565b905085602001518111156110555761105586610f42836002611538565b85518051838201602001916000918085111561106f578482525b505050602086015b602086106110af578051825261108e6020836116bd565b915061109b6020826116bd565b90506110a86020876116aa565b9550611077565b5181517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff60208890036101000a0190811690199190911617905250849150509392505050565b81516111018383610c9b565b50610dd98382610ed9565b6040518060400160405280611134604051806040016040528060608152602001600081525090565b8152602001600081525090565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156111b7576111b7611141565b604052919050565b600067ffffffffffffffff8311156111d9576111d9611141565b61120a60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f86011601611170565b905082815283838301111561121e57600080fd5b828260208301376000602084830101529392505050565b600082601f83011261124657600080fd5b610efa838335602085016111bf565b60008060006060848603121561126a57600080fd5b83359250602084013567ffffffffffffffff8082111561128957600080fd5b61129587838801611235565b935060408601359150808211156112ab57600080fd5b506112b886828701611235565b9150509250925092565b60008083601f8401126112d457600080fd5b50813567ffffffffffffffff8111156112ec57600080fd5b60208301915083602082850101111561130457600080fd5b9250929050565b803567ffffffffffffffff8116811461132357600080fd5b919050565b60008060008060008060008060a0898b03121561134457600080fd5b883567ffffffffffffffff8082111561135c57600080fd5b6113688c838d016112c2565b909a50985060208b013591508082111561138157600080fd5b61138d8c838d016112c2565b909850965060408b01359150808211156113a657600080fd5b818b0191508b601f8301126113ba57600080fd5b8135818111156113c957600080fd5b8c60208260051b85010111156113de57600080fd5b6020830196508095505050506113f660608a0161130b565b9150608089013590509295985092959890939650565b60006020828403121561141e57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff81168114610efa57600080fd5b600067ffffffffffffffff8084111561145d5761145d611141565b8360051b602061146e818301611170565b86815291850191818101903684111561148657600080fd5b865b848110156114ce578035868111156114a05760008081fd5b880136601f8201126114b25760008081fd5b6114c03682358784016111bf565b845250918301918301611488565b50979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610d0f57610d0f611509565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115af576115af611509565b5060010190565b67ffffffffffffffff861681526000602060a08184015286518060a085015260005b818110156115f45788810183015185820160c0015282016115d8565b50600060c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505061163d604083018661ffff169052565b63ffffffff939093166060820152608001529392505050565b60006020828403121561166857600080fd5b5051919050565b6000826116a5577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500690565b81810381811115610d0f57610d0f611509565b80820180821115610d0f57610d0f611509565b600181815b8085111561172957817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561170f5761170f611509565b8085161561171c57918102915b93841c93908002906116d5565b509250929050565b60008261174057506001610d0f565b8161174d57506000610d0f565b8160018114611763576002811461176d57611789565b6001915050610d0f565b60ff84111561177e5761177e611509565b50506001821b610d0f565b5060208310610133831016604e8410600b84101617156117ac575081810a610d0f565b6117b683836116d0565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156117e8576117e8611509565b029392505050565b6000610efa838361173156fea164736f6c6343000813000a",
}

var FunctionsClientExampleABI = FunctionsClientExampleMetaData.ABI

var FunctionsClientExampleBin = FunctionsClientExampleMetaData.Bin

func DeployFunctionsClientExample(auth *bind.TransactOpts, backend bind.ContractBackend, router common.Address) (common.Address, *types.Transaction, *FunctionsClientExample, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FunctionsClientExampleBin), backend, router)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FunctionsClientExample{FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

type FunctionsClientExample struct {
	address common.Address
	abi     abi.ABI
	FunctionsClientExampleCaller
	FunctionsClientExampleTransactor
	FunctionsClientExampleFilterer
}

type FunctionsClientExampleCaller struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleTransactor struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleFilterer struct {
	contract *bind.BoundContract
}

type FunctionsClientExampleSession struct {
	Contract     *FunctionsClientExample
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleCallerSession struct {
	Contract *FunctionsClientExampleCaller
	CallOpts bind.CallOpts
}

type FunctionsClientExampleTransactorSession struct {
	Contract     *FunctionsClientExampleTransactor
	TransactOpts bind.TransactOpts
}

type FunctionsClientExampleRaw struct {
	Contract *FunctionsClientExample
}

type FunctionsClientExampleCallerRaw struct {
	Contract *FunctionsClientExampleCaller
}

type FunctionsClientExampleTransactorRaw struct {
	Contract *FunctionsClientExampleTransactor
}

func NewFunctionsClientExample(address common.Address, backend bind.ContractBackend) (*FunctionsClientExample, error) {
	abi, err := abi.JSON(strings.NewReader(FunctionsClientExampleABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindFunctionsClientExample(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExample{address: address, abi: abi, FunctionsClientExampleCaller: FunctionsClientExampleCaller{contract: contract}, FunctionsClientExampleTransactor: FunctionsClientExampleTransactor{contract: contract}, FunctionsClientExampleFilterer: FunctionsClientExampleFilterer{contract: contract}}, nil
}

func NewFunctionsClientExampleCaller(address common.Address, caller bind.ContractCaller) (*FunctionsClientExampleCaller, error) {
	contract, err := bindFunctionsClientExample(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleCaller{contract: contract}, nil
}

func NewFunctionsClientExampleTransactor(address common.Address, transactor bind.ContractTransactor) (*FunctionsClientExampleTransactor, error) {
	contract, err := bindFunctionsClientExample(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleTransactor{contract: contract}, nil
}

func NewFunctionsClientExampleFilterer(address common.Address, filterer bind.ContractFilterer) (*FunctionsClientExampleFilterer, error) {
	contract, err := bindFunctionsClientExample(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleFilterer{contract: contract}, nil
}

func bindFunctionsClientExample(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FunctionsClientExampleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.FunctionsClientExampleCaller.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.FunctionsClientExampleTransactor.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FunctionsClientExample.Contract.contract.Call(opts, result, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transfer(opts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.contract.Transact(opts, method, params...)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "MAX_CALLBACK_GAS")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) MAXCALLBACKGAS() (uint32, error) {
	return _FunctionsClientExample.Contract.MAXCALLBACKGAS(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) Owner() (common.Address, error) {
	return _FunctionsClientExample.Contract.Owner(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastError(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastError")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastError() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastError(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastErrorLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastErrorLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastErrorLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastErrorLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastRequestId(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastRequestId")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastRequestId() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastRequestId(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastResponse(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastResponse")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastResponse() ([32]byte, error) {
	return _FunctionsClientExample.Contract.SLastResponse(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCaller) SLastResponseLength(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FunctionsClientExample.contract.Call(opts, &out, "s_lastResponseLength")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

func (_FunctionsClientExample *FunctionsClientExampleSession) SLastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleCallerSession) SLastResponseLength() (uint32, error) {
	return _FunctionsClientExample.Contract.SLastResponseLength(&_FunctionsClientExample.CallOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "acceptOwnership")
}

func (_FunctionsClientExample *FunctionsClientExampleSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.AcceptOwnership(&_FunctionsClientExample.TransactOpts)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "handleOracleFulfillment", requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) HandleOracleFulfillment(requestId [32]byte, response []byte, err []byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.HandleOracleFulfillment(&_FunctionsClientExample.TransactOpts, requestId, response, err)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "sendRequest", source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) SendRequest(source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.SendRequest(&_FunctionsClientExample.TransactOpts, source, encryptedSecretsReferences, args, subscriptionId, jobId)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactor) TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.contract.Transact(opts, "transferOwnership", to)
}

func (_FunctionsClientExample *FunctionsClientExampleSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

func (_FunctionsClientExample *FunctionsClientExampleTransactorSession) TransferOwnership(to common.Address) (*types.Transaction, error) {
	return _FunctionsClientExample.Contract.TransferOwnership(&_FunctionsClientExample.TransactOpts, to)
}

type FunctionsClientExampleOwnershipTransferRequestedIterator struct {
	Event *FunctionsClientExampleOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
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
		it.Event = new(FunctionsClientExampleOwnershipTransferRequested)
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

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferRequestedIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferRequested)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error) {
	event := new(FunctionsClientExampleOwnershipTransferRequested)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleOwnershipTransferredIterator struct {
	Event *FunctionsClientExampleOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleOwnershipTransferred)
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
		it.Event = new(FunctionsClientExampleOwnershipTransferred)
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

func (it *FunctionsClientExampleOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleOwnershipTransferredIterator{contract: _FunctionsClientExample.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleOwnershipTransferred)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error) {
	event := new(FunctionsClientExampleOwnershipTransferred)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestFulfilledIterator struct {
	Event *FunctionsClientExampleRequestFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestFulfilled)
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
		it.Event = new(FunctionsClientExampleRequestFulfilled)
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

func (it *FunctionsClientExampleRequestFulfilledIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestFulfilled struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestFulfilledIterator{contract: _FunctionsClientExample.contract, event: "RequestFulfilled", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestFulfilled", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestFulfilled)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error) {
	event := new(FunctionsClientExampleRequestFulfilled)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type FunctionsClientExampleRequestSentIterator struct {
	Event *FunctionsClientExampleRequestSent

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *FunctionsClientExampleRequestSentIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FunctionsClientExampleRequestSent)
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
		it.Event = new(FunctionsClientExampleRequestSent)
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

func (it *FunctionsClientExampleRequestSentIterator) Error() error {
	return it.fail
}

func (it *FunctionsClientExampleRequestSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type FunctionsClientExampleRequestSent struct {
	Id  [32]byte
	Raw types.Log
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.FilterLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return &FunctionsClientExampleRequestSentIterator{contract: _FunctionsClientExample.contract, event: "RequestSent", logs: logs, sub: sub}, nil
}

func (_FunctionsClientExample *FunctionsClientExampleFilterer) WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _FunctionsClientExample.contract.WatchLogs(opts, "RequestSent", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(FunctionsClientExampleRequestSent)
				if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
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

func (_FunctionsClientExample *FunctionsClientExampleFilterer) ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error) {
	event := new(FunctionsClientExampleRequestSent)
	if err := _FunctionsClientExample.contract.UnpackLog(event, "RequestSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_FunctionsClientExample *FunctionsClientExample) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _FunctionsClientExample.abi.Events["OwnershipTransferRequested"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferRequested(log)
	case _FunctionsClientExample.abi.Events["OwnershipTransferred"].ID:
		return _FunctionsClientExample.ParseOwnershipTransferred(log)
	case _FunctionsClientExample.abi.Events["RequestFulfilled"].ID:
		return _FunctionsClientExample.ParseRequestFulfilled(log)
	case _FunctionsClientExample.abi.Events["RequestSent"].ID:
		return _FunctionsClientExample.ParseRequestSent(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (FunctionsClientExampleOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (FunctionsClientExampleOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (FunctionsClientExampleRequestFulfilled) Topic() common.Hash {
	return common.HexToHash("0x85e1543bf2f84fe80c6badbce3648c8539ad1df4d2b3d822938ca0538be727e6")
}

func (FunctionsClientExampleRequestSent) Topic() common.Hash {
	return common.HexToHash("0x1131472297a800fee664d1d89cfa8f7676ff07189ecc53f80bbb5f4969099db8")
}

func (_FunctionsClientExample *FunctionsClientExample) Address() common.Address {
	return _FunctionsClientExample.address
}

type FunctionsClientExampleInterface interface {
	MAXCALLBACKGAS(opts *bind.CallOpts) (uint32, error)

	Owner(opts *bind.CallOpts) (common.Address, error)

	SLastError(opts *bind.CallOpts) ([32]byte, error)

	SLastErrorLength(opts *bind.CallOpts) (uint32, error)

	SLastRequestId(opts *bind.CallOpts) ([32]byte, error)

	SLastResponse(opts *bind.CallOpts) ([32]byte, error)

	SLastResponseLength(opts *bind.CallOpts) (uint32, error)

	AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error)

	HandleOracleFulfillment(opts *bind.TransactOpts, requestId [32]byte, response []byte, err []byte) (*types.Transaction, error)

	SendRequest(opts *bind.TransactOpts, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) (*types.Transaction, error)

	TransferOwnership(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*FunctionsClientExampleOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FunctionsClientExampleOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*FunctionsClientExampleOwnershipTransferred, error)

	FilterRequestFulfilled(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestFulfilledIterator, error)

	WatchRequestFulfilled(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestFulfilled, id [][32]byte) (event.Subscription, error)

	ParseRequestFulfilled(log types.Log) (*FunctionsClientExampleRequestFulfilled, error)

	FilterRequestSent(opts *bind.FilterOpts, id [][32]byte) (*FunctionsClientExampleRequestSentIterator, error)

	WatchRequestSent(opts *bind.WatchOpts, sink chan<- *FunctionsClientExampleRequestSent, id [][32]byte) (event.Subscription, error)

	ParseRequestSent(log types.Log) (*FunctionsClientExampleRequestSent, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
