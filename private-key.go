package eth

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewPrivateKey() PrivateKey {
	privateKeyECDSA, err := crypto.GenerateKey()
	if err != nil {
		log.Println(err)
		return nil
	}
	return &privateKey{ecdsaPrivateKey: privateKeyECDSA}
}

func NewPrivateKeyFromHex(hex string) PrivateKey {
	if len(hex) >= 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	privateKeyECDSA, err := crypto.HexToECDSA(hex)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &privateKey{ecdsaPrivateKey: privateKeyECDSA}
}

type PrivateKey interface {
	Hex() string
	PublicKey() PublicKey
	RawPrivateKey() *ecdsa.PrivateKey
}

type privateKey struct {
	ecdsaPrivateKey *ecdsa.PrivateKey
}

func (p *privateKey) Hex() string {
	return hexutil.Encode(crypto.FromECDSA(p.ecdsaPrivateKey))[2:]
}

func (p *privateKey) PublicKey() PublicKey {
	return &publicKey{
		ecdsaPublicKey: &p.ecdsaPrivateKey.PublicKey,
	}
}

func (p *privateKey) RawPrivateKey() *ecdsa.PrivateKey {
	return p.ecdsaPrivateKey
}
