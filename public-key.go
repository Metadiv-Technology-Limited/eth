package eth

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func NewPublicKeyFromHex(hex string) PublicKey {
	publicKeyECDSA, err := crypto.HexToECDSA(hex)
	if err != nil {
		panic(err)
	}
	return &publicKey{
		ecdsaPublicKey: &publicKeyECDSA.PublicKey,
	}

}

type PublicKey interface {
	Address() string
	RawPublicKey() *ecdsa.PublicKey
}

type publicKey struct {
	ecdsaPublicKey *ecdsa.PublicKey
}

func (p *publicKey) RawPublicKey() *ecdsa.PublicKey {
	return p.ecdsaPublicKey
}

func (p *publicKey) Address() string {
	return crypto.PubkeyToAddress(*p.ecdsaPublicKey).Hex()
}
