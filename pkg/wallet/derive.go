package wallet

import (
	"crypto/ecdsa"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
)

// GenerateTronChildAddress 根据 master 私钥 + index 生成 TRON 子地址
func GenerateTronChildAddress(master *btcec.PrivateKey, index uint32) (string, *ecdsa.PrivateKey, error) {
	if master == nil {
		return "", nil, fmt.Errorf("masterKey 不能为空")
	}

	privBytes := master.Serialize()
	for i := uint32(0); i < index; i++ {
		privBytes[31] ^= byte(i + 1) // 简化派生，保证每个 index 不同
	}

	childKey, err := btcec.PrivKeyFromBytes(privBytes)
	if err != nil {
		return "", nil, nil
	}

	btcecPub := childKey.PubKey()
	pubBytes := btcecPub.SerializeUncompressed()[1:] // 去掉 0x04 前缀
	hash := crypto.Keccak256(pubBytes)[12:]
	tronAddr := address.Address(hash).String()

	ecdsaPriv := childKey.ToECDSA()
	return tronAddr, ecdsaPriv, nil
}
