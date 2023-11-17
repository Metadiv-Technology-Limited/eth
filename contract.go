package eth

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewContract(abi string, bin string, address string) Contract {
	return &contract{abi: abi, bin: bin, address: address}
}

type Contract interface{}

type contract struct {
	abi     string
	bin     string
	address string
}

func (c *contract) GetConstructorMethod() *Abi {
	for _, method := range c.abiJson() {
		if method.Type == "constructor" {
			return &method
		}
	}
	return nil
}

func (c *contract) GetMethodByName(name string) *Abi {
	for _, method := range c.abiJson() {
		if method.Name == name {
			return &method
		}
	}
	return nil
}

func (c *contract) GetAllReadMethods() []Abi {
	var methods []Abi
	for _, method := range c.abiJson() {
		if method.Type == "function" && method.StateMutability == "view" {
			methods = append(methods, method)
		}
	}
	return methods
}

func (c *contract) GetAllWriteMethods() []Abi {
	var methods []Abi
	for _, method := range c.abiJson() {
		if method.Type == "function" && method.StateMutability != "view" {
			methods = append(methods, method)
		}
	}
	return methods
}

func (c *contract) Deploy(client Client, opts TransactOpts, args ...any) (*DeployContractResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	address, tx, _, err := bind.DeployContract(
		c.setupTransactOpts(client, opts, ctx), c.parseAbi(), c.parseBin(), client.RawEthClient(), args...)
	if err != nil {
		return nil, err
	}

	return &DeployContractResult{
		Address: address.Hex(),
		TxHash:  tx.Hash().Hex(),
		Nonce:   tx.Nonce(),
	}, nil
}

func (c *contract) Read(client Client, method string, args ...interface{}) ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	contract := bind.NewBoundContract(common.HexToAddress(c.address), c.parseAbi(), client.RawEthClient(), client.RawEthClient(), client.RawEthClient())

	var result []interface{}
	err := contract.Call(&bind.CallOpts{Context: ctx}, &result, method, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *contract) Write(client Client, opts TransactOpts, method string, args ...interface{}) (*WriteContractResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	parsed, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		log.Fatal(err)
	}
	contract := bind.NewBoundContract(common.HexToAddress(c.address), parsed, client.RawEthClient(), client.RawEthClient(), client.RawEthClient())

	tx, err := contract.Transact(c.setupTransactOpts(client, opts, ctx), method, args...)
	if err != nil {
		return nil, err
	}

	return &WriteContractResult{
		TxHash: tx.Hash().Hex(),
		Nonce:  tx.Nonce(),
	}, nil
}

func (c *contract) abiJson() []Abi {
	abiJsonList := make([]Abi, 0)
	c.abi = strings.ReplaceAll(c.abi, "\\\"", "\"")
	err := json.Unmarshal([]byte(c.abi), &abiJsonList)
	if err != nil {
		panic(err)
	}
	return abiJsonList
}

func (c *contract) parseAbi() abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		panic(err)
	}
	return parsed
}

func (c *contract) parseBin() []byte {
	return common.FromHex(c.bin)
}

func (c *contract) setupTransactOpts(client Client, opts TransactOpts, ctx context.Context) *bind.TransactOpts {
	clientID, err := client.ChainID()
	if err != nil {
		log.Println(err)
		return nil
	}

	opt, err := bind.NewKeyedTransactorWithChainID(client.PrivateKey().RawPrivateKey(), big.NewInt(int64(clientID)))
	if err != nil {
		log.Println(err)
		return nil
	}
	opt.From = crypto.PubkeyToAddress(*client.PrivateKey().PublicKey().RawPublicKey())
	opt.Value = big.NewInt(int64(opts.Value))
	opt.GasLimit = uint64(opts.GasLimit)
	opt.GasPrice = big.NewInt(int64(opts.GasPrice))

	return opt
}
