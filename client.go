package eth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(rpcURL string, privateKey PrivateKey) Client {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &client{ethClient: ethClient, privateKey: privateKey}
}

func NewClientWithoutPrivateKey(rpcURL string) Client {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &client{ethClient: ethClient}
}

type Client interface {
	RawEthClient() *ethclient.Client
	PrivateKey() PrivateKey
	Transaction(txHash string) (tx *types.Transaction, isPending bool, err error)
	BalanceOf(address string) (uint64, error)
	ChainID() (uint64, error)
	SuggestedGasPrice() (uint64, error)
	WaitForMined(txHash string) error
}

type client struct {
	privateKey PrivateKey
	ethClient  *ethclient.Client
}

func (c *client) RawEthClient() *ethclient.Client {
	return c.ethClient
}

func (c *client) PrivateKey() PrivateKey {
	return c.privateKey
}

func (c *client) Transaction(txHash string) (tx *types.Transaction, isPending bool, err error) {
	tx, isPending, err = c.ethClient.TransactionByHash(context.Background(), common.HexToHash(txHash))
	return
}

func (c *client) BalanceOf(address string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	balance, err := c.RawEthClient().BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (c *client) ChainID() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	chainID, err := c.RawEthClient().NetworkID(ctx)
	if err != nil {
		return 0, err
	}
	return chainID.Uint64(), nil
}

func (c *client) SuggestedGasPrice() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	gasPrice, err := c.RawEthClient().SuggestGasPrice(ctx)
	if err != nil {
		return 0, err
	}
	return gasPrice.Uint64(), nil
}

func (c *client) WaitForMined(txHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tx, _, err := c.Transaction(txHash)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(ctx, c.RawEthClient(), tx)
	if err != nil {
		return err
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return errors.New("transaction failed: " + tx.Hash().Hex())
	}
	return nil
}
