package service

import (
	"fmt"
	"log"
	"math/big"
	"time"
	"tron-scan/config"

	_ "time"

	"tron-scan/internal/httpserver/handler/response"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	trxgrpcs "tron-scan/pkg/trx/grpcs"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	_ "google.golang.org/grpc"
)

// ---------------------------
// 助记词常量（程序内固定）
// ---------------------------
const MASTER_MNEMONIC = "fame combine setup people again basket enlist drive live dad share debate"

type (
	Service interface {
		// 随机生成 Tron 钱包
		CreateUserWallet() (*WalletResult, error)
		// 固定助记词生成主钱包（可派生子地址）
		CreateUserMainWallet() (*WalletResult, error)
		// 根据索引派生子地址
		DeriveUserSubWallet(index uint32) (*WalletResult, error)
		// 查询
		QueryTronAccountAssets(tronAddr string) (*response.TronAssetResult, error)
	}
)

// WalletResult 返回给前端或业务层的钱包信息
type WalletResult struct {
	Mnemonic    string // 助记词
	PrivateKey  string // 私钥 hex
	TronAddress string // Tron Base58 地址
}

// WalletService 结构体
type WalletService struct{}

// ---------------------------
// 1️⃣ 随机生成 Tron 钱包
// ---------------------------
func (s *WalletService) CreateUserWallet() (*WalletResult, error) {
	entropy, err := bip39.NewEntropy(128) // 128 位 entropy -> 12 个助记词
	if err != nil {
		return nil, fmt.Errorf("生成助记词 entropy 失败: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("生成助记词失败: %w", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	privKey, err := crypto.ToECDSA(seed[:32])
	if err != nil {
		return nil, fmt.Errorf("生成私钥失败: %w", err)
	}

	// 使用 gotron-sdk 生成 Tron 地址
	tronAddr := address.PubkeyToAddress(privKey.PublicKey).String()

	result := &WalletResult{
		Mnemonic:    mnemonic,
		PrivateKey:  fmt.Sprintf("%x", crypto.FromECDSA(privKey)),
		TronAddress: tronAddr,
	}

	log.Printf("随机 Tron 钱包创建成功: %+v", result)
	return result, nil
}

// ---------------------------
// 2️⃣ 固定助记词生成主钱包
// ---------------------------
func (s *WalletService) CreateUserMainWallet() (*WalletResult, error) {
	return DeriveTronWallet(MASTER_MNEMONIC, 0)
}

// ---------------------------
// 3️⃣ 根据索引派生子地址（BIP44）
// ---------------------------
func (s *WalletService) DeriveUserSubWallet(index uint32) (*WalletResult, error) {
	return DeriveTronWallet(MASTER_MNEMONIC, index)
}

// ---------------------------
// 公共函数: 根据助记词和 index 派生 Tron 钱包
// ---------------------------
func DeriveTronWallet(mnemonic string, index uint32) (*WalletResult, error) {
	// 1️⃣ 助记词 -> seed
	seed := bip39.NewSeed(mnemonic, "")

	// 2️⃣ master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("生成主私钥失败: %w", err)
	}

	// 3️⃣ BIP44 路径 m/44'/195'/0'/0/index
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 195) // TRON coin_type=195
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	addrKey, _ := change.NewChildKey(index)

	// 4️⃣ 私钥 -> 公钥 -> Tron 地址
	privKey, err := crypto.ToECDSA(addrKey.Key)
	if err != nil {
		return nil, fmt.Errorf("生成 ECDSA 私钥失败: %w", err)
	}
	tronAddr := address.PubkeyToAddress(privKey.PublicKey).String()

	return &WalletResult{
		Mnemonic:    "",
		PrivateKey:  "",
		TronAddress: tronAddr,
	}, nil
}

func (s *WalletService) QueryTronAccountAssets(tronAddr string) (*response.TronAssetResult, error) {
	// ===== 1️⃣ 获取客户端 =====
	client, err := trxgrpcs.GetClient()
	if err != nil {
		return nil, fmt.Errorf("create tron client failed: %w", err)
	}
	// ===== 2️⃣ 查询 TRX =====
	amountTRX, err := client.GetTrxBalance(tronAddr)
	if err != nil || amountTRX == nil {
		// 防止空指针
		amountTRX = big.NewInt(0)
	}
	// TRX 单位转换（Sun → TRX）
	trxSun := amountTRX.Int64()
	trxFloat := new(big.Float).
		Quo(new(big.Float).SetInt(amountTRX), big.NewFloat(1e6)).
		Text('f', 6)
	// ===== 3️⃣ 查询 TRC20（例如 USDT）=====
	trc20Map := make(map[string]string)

	// Nile测试网USDT合约地址
	//usdtContract := "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"
	usdtContract := config.GlobalConfig.Tron.USDTContract

	balanceInt, err := client.GetTrc20Balance(tronAddr, usdtContract)
	if err == nil && balanceInt != nil {

		// USDT 6位小数
		usdtFloat := new(big.Float).
			Quo(new(big.Float).SetInt(balanceInt), big.NewFloat(1e6)).
			Text('f', 6)

		trc20Map["USDT"] = usdtFloat
	} else {
		trc20Map["USDT"] = "0"
	}
	// 查询资源
	resource, errRes := client.GetAccountResource(tronAddr)
	var bandwidth, bandwidthUsed, energy, energyUsed int64

	if errRes == nil && resource != nil {
		// 免费带宽
		bandwidth = resource.FreeNetLimit
		// 已用带宽
		bandwidthUsed = resource.FreeNetUsed
		// 能量（TRC20用）
		energy = resource.EnergyLimit
		// 已用能量
		energyUsed = resource.EnergyUsed
	}
	// ===== 5️⃣ 返回结果 =====
	result := &response.TronAssetResult{
		Address:       tronAddr,      // 钱包地址
		TRX:           trxSun,        // TRX（Sun）
		TRXFloat:      trxFloat,      // TRX（可读）
		Bandwidth:     bandwidth,     // 带宽总量
		BandwidthUsed: bandwidthUsed, // 已用带宽
		Energy:        energy,        // 能量总量
		EnergyUsed:    energyUsed,    // 已用能量
		TRC20:         trc20Map,      // TRC20资产
	}
	return result, nil
}

func (s *WalletService) GetTRC20TxHistory(userAddr string) ([]*response.TRC20Transfer, error) {
	// ===== 1️⃣ 获取客户端 =====
	c, err := trxgrpcs.GetClient()
	if err != nil {
		return nil, fmt.Errorf("create tron client failed: %w", err)
	}
	var allData []*response.TRC20Transfer
	currentFingerprint := ""

	for {
		// 调用你刚才写的函数
		result, err := c.GetTRC20TxHistory(userAddr, currentFingerprint)
		if err != nil {
			return nil, err
		}

		// 将当前页数据加入总列表
		allData = append(allData, result.Transfers...)

		// 打印进度 (可选)
		fmt.Printf("已拉取 %d 条...\n", len(allData))

		// 核心逻辑：如果没有指纹了，说明拉完了
		if !result.HasMore {
			break
		}

		// 更新指纹，进入下一轮循环
		currentFingerprint = result.Fingerprint

		// 【重要】频率控制：如果你没有 API Key 或者 Key 级别较低，
		// 建议增加一小段休眠，防止被 TronGrid 封禁 IP
		time.Sleep(200 * time.Millisecond)
	}

	return allData, nil
}

func (s *WalletService) TransferTRX(fromAddr string, toAddr string, amount int64) (string, error) {
	// ===== 1️⃣ 获取客户端 =====
	c, err := trxgrpcs.GetClient()
	if err != nil {
		return "", fmt.Errorf("create tron client failed: %w", err)
	}
	privateKeyHex := "691e9dadd84a55067ac312678c1737457e5f260c5e1f9de81af23ccb5cd741c9"
	result, err := c.TransferTRX(fromAddr, toAddr, privateKeyHex, amount)
	if err != nil {
		return "", fmt.Errorf("转账TRX失败: %w", err)
	}
	return result, nil
}

func (s *WalletService) TransferTRC20(fromAddr string, toAddr string, amount int64) (string, error) {
	// ===== 1️⃣ 获取客户端 =====
	c, err := trxgrpcs.GetClient()
	if err != nil {
		return "", fmt.Errorf("create tron client failed: %w", err)
	}
	privateKeyHex := "691e9dadd84a55067ac312678c1737457e5f260c5e1f9de81af23ccb5cd741c9"
	contractAddr := "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"
	//4. 设置手续费上限 (建议 100 TRX，即 100,000,000 Sun)
	// 因为波场现在转账 USDT 比较贵，设置低了会 OUT_OF_ENERGY
	var feeLimit int64 = 100 * 1000000
	result, err := c.TransferTRC20(fromAddr, toAddr, contractAddr, privateKeyHex, big.NewInt(amount), feeLimit)
	if err != nil {
		return "", fmt.Errorf("转账TRX失败: %w", err)
	}
	return result, nil
}
