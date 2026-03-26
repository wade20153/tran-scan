package tron

import (
	"encoding/hex"
	
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/tyler-smith/go-bip39"
)

type TronAccount struct {
	Mnemonic   string // 助记词（12词）
	PrivateKey string // hex 私钥
	Address    string // Base58 地址（T开头）
}

func CreateTronAccount() (*TronAccount, error) {
	// 1️⃣ 生成 128bit 随机熵（12个助记词）
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}

	// 2️⃣ 生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	// 3️⃣ 助记词 → Seed
	seed := bip39.NewSeed(mnemonic, "")

	// 4️⃣ 生成主私钥（⚠️ TRON/ETH 都是 secp256k1）
	masterKey, _ := btcec.PrivKeyFromBytes(seed[:32])

	// 5️⃣ 私钥 hex
	privHex := hex.EncodeToString(masterKey.Serialize())

	// 6️⃣ 公钥 → TRON 地址
	pubBytes := masterKey.PubKey().SerializeUncompressed()[1:]
	hash := crypto.Keccak256(pubBytes)[12:]
	tronAddr := address.Address(hash).String()

	return &TronAccount{
		Mnemonic:   mnemonic,
		PrivateKey: privHex,
		Address:    tronAddr,
	}, nil
}
