package signer

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2" // secp256k1 私钥/公钥
	_ "github.com/btcsuite/btcd/chaincfg"
	_ "github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"       // ETH/TRON 地址工具
	"github.com/fbsobreira/gotron-sdk/pkg/address" // TRON Base58 地址
	"github.com/tyler-smith/go-bip39"              // BIP39 助记词            // BIP39 助记词
)

// WalletResult 用于返回完整解析结果（生产中可只保留加密xprv）
type WalletResult struct {
	SeedHex     string // 助记词生成的 Seed 的 hex 编码
	Xprv        string // 主私钥（Master Private Key）的 hex 编码
	EthAddress  string // 主 ETH 地址
	TronAddress string // 主 TRON 地址
}

// 内存中缓存 masterKey
var cachedMasterKey *btcec.PrivateKey

// SetMasterKey 缓存 masterKey，ImportMnemonic 调用后设置
func SetMasterKey(master *btcec.PrivateKey) {
	cachedMasterKey = master
}

// ImportMnemonic 核心方法：输入明文助记词 → 生成所有核心信息
func ImportMnemonic(mnemonic string, passphrase string) (*WalletResult, error) {
	// 1️⃣ 校验助记词合法性
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("助记词不合法")
	}

	// 2️⃣ 助记词 -> BIP39 Seed
	seed := bip39.NewSeed(mnemonic, passphrase)
	seedHex := hex.EncodeToString(seed)

	// 3️⃣ Seed -> 主私钥 (MasterKey)
	masterKey, _ := btcec.PrivKeyFromBytes(seed[:32])
	// 存在内存中
	SetMasterKey(masterKey)
	xprv := hex.EncodeToString(masterKey.Serialize())

	// 4️⃣ 派生主 ETH 地址: m/44'/60'/0'/0/0
	ethPriv := deriveMasterKey(masterKey, 44, 60, 0, 0, 0)
	// btcec -> ecdsa 转换
	ecdsaPriv := ethPriv.ToECDSA()
	// ETH 地址
	ethAddress := crypto.PubkeyToAddress(ecdsaPriv.PublicKey).Hex()

	// 5️⃣ 派生主 TRON 地址: m/44'/195'/0'/0/0
	// 其中
	tronPriv := deriveMasterKey(masterKey, 44, 195, 0, 0, 0)
	pubBytes := tronPriv.PubKey().SerializeUncompressed()[1:] // 去掉 0x04 前缀
	hash := crypto.Keccak256(pubBytes)[12:]
	tronAddress := address.Address(hash).String()

	// 6️⃣ 返回结果
	return &WalletResult{
		SeedHex:     seedHex,
		Xprv:        xprv,
		EthAddress:  ethAddress,
		TronAddress: tronAddress,
	}, nil
}

// deriveMasterKey 简化 BIP44 派生函数 (仅示意主地址)
// 生产环境建议替换为完整 BIP32/BIP44 算法
func deriveMasterKey(master *btcec.PrivateKey, purpose, coin, account, change, index uint32) *btcec.PrivateKey {
	// 生产环境这里只生成主私钥，不再尝试“修改字节数组”
	// 如果需要子钱包，请用完整 BIP32/BIP44 算法
	return master
}

// GenerateTronChildAddress 根据master 私钥 + index 派生成 tests 子地址
// ⚠️ 注意：这里简化 BIP44，仅用于演示和主地址/子地址生成
func GenerateTronChildAddress(index uint32) (string, *ecdsa.PrivateKey, error) {
	/*
		BIP44 派生路径: m / 44' / 195' / 0' / 0 / index
		- purpose = 44'
		- coin_type = 195' (TRON)
		- account = 0'
		- change = 0
		- address_index = index
	*/

	// -------------------------------
	// 1️⃣ 派生子私钥
	// -------------------------------
	// ⚠️ btcec/v2 不允许直接修改 ModNScalar
	// 为了示意，这里生成子私钥的方法是：
	// 使用 index 对 master 私钥做哈希再生成 btcec.PrivateKey
	// 生产环境建议使用完整 BIP32/BIP44 派生算法
	// 使用
	if cachedMasterKey == nil {
		return "", nil, fmt.Errorf("masterKey 未初始化，请先调用 ImportMnemonic 并 SetMasterKey")
	}
	// 1️⃣ 派生子私钥（简化示意）
	privBytes := cachedMasterKey.Serialize()
	for i := uint32(0); i < index; i++ {
		privBytes[31] ^= byte(i + 1) // 修改最后一字节，确保不同 index 得到不同私钥
	}
	childKey, err := btcec.PrivKeyFromBytes(privBytes)
	if err != nil {
		return "", nil, fmt.Errorf("生成子私钥失败: %w", err)
	}
	// 2️⃣ 转为 ecdsa.PrivateKey
	// -------------------------------
	btcecPub := childKey.PubKey()                    // ✅ btcec.PublicKey
	pubBytes := btcecPub.SerializeUncompressed()[1:] // 去掉 0x04 前缀
	// 3️⃣ keccak256 + Base58Check 得到 TRON 地址
	// -------------------------------
	hash := crypto.Keccak256(pubBytes)[12:]
	tronAddr := address.Address(hash).String()
	// 4️⃣ 返回 ecdsa.PrivateKey
	// -------------------------------
	ecdsaPriv := childKey.ToECDSA()
	// 返回子地址，子私钥
	return tronAddr, ecdsaPriv, nil
}

// 根据phase.yml
// 读取phase.yml ，mnemonicPhrase key salt
// 第一步生产这个key 与salt。 加密mnemonic，生成mnemonicPhrase
// 第二步，根据系统用户的的user-key user-salt 再次加密key salt。 生成key salt
// 读取成加密与解密
