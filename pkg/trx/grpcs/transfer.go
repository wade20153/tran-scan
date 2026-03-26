package grpcs

import (
	"encoding/hex"
	"fmt"

	_ "github.com/fbsobreira/gotron-sdk/pkg/proto/core"

	"github.com/btcsuite/btcd/btcec/v2"
	_ "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	_ "github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

// TransferTRX 转账 TRX（单位：SUN）
func (c *Client) TransferTRX(
	fromAddr string,
	toAddr string,
	privateKeyHex string,
	amountSun int64,
) (string, error) {
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return "error", err
	}
	// 1️⃣ 创建交易 1 TRX = 1,000,000 SUN
	amountSun = int64(amountSun * 1_000_000)
	txExt, err := grpcCli.Transfer(fromAddr, toAddr, amountSun)
	if err != nil {
		return "", fmt.Errorf("create transaction failed: %w", err)
	}

	// 2️⃣ 将私钥 hex 转 bytes
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	// 3️⃣ 解析私钥
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// 4️⃣ 使用 transaction 包签名
	signedTx, err := transaction.SignTransaction(txExt.Transaction, privKey)
	if err != nil {
		return "", fmt.Errorf("sign transaction failed: %w", err)
	}

	// 将签名放回 TransactionExtention
	txExt.Transaction = signedTx

	// 5️⃣ 广播交易
	result, err := grpcCli.Broadcast(txExt.Transaction)
	if err != nil {
		return "", fmt.Errorf("broadcast transaction failed: %w", err)
	}
	if !result.Result {
		return "", fmt.Errorf("broadcast trx failed")
	}

	// 6️⃣ 获取 txid
	txidBytes := txExt.GetTxid() // ✅ TransactionExtention 自带方法
	return hex.EncodeToString(txidBytes), nil
}
