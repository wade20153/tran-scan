package grpcs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	//"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa" // 👈 必须导入这个
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"google.golang.org/protobuf/proto"
)

// TransferTRC20 使用 TRON 客户端转账 TRC20 Token
// 参数说明：
//
//		fromAddr: 发送方钱包地址
//		toAddr: 接收方钱包地址
//		contractAddr: TRC20 合约地址
//		privateKeyHex: 发送方私钥（16 进制字符串）
//		amount: 转账金额（big.Int，单位为最小单位，例如 Sun）
//		feeLimit: 最大能支付的交易手续费（sun）
//	 返回值：
//		string: 返回交易 ID（txid）
//		error: 如果发生错误，返回具体错误信息
func (c *Client) TransferTRC20(
	from string,
	to string,
	contractAddr string,
	privateKeyHex string,
	amount *big.Int,
	feeLimit int64,
) (string, error) {
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return "error", err
	}

	fromAddr, _ := address.Base58ToAddress(from)
	toAddr, _ := address.Base58ToAddress(to)
	data := fmt.Sprintf(`[{"address": "%s"}, {"uint256": "0x%x"}]`, toAddr.String(), amount)

	// 4️⃣ 调用 TriggerContract
	txExt, err := grpcCli.TriggerContract(
		fromAddr.String(),
		contractAddr,
		"transfer(address,uint256)",
		data,
		feeLimit,
		0, "", 0,
	)
	if err != nil {
		return "", fmt.Errorf("trigger TRC20 contract failed: %w", err)
	}

	// 3️⃣ 准备私钥
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// 4️⃣ 修正：计算交易 RawData 的哈希
	rawDataBytes, err := proto.Marshal(txExt.Transaction.RawData)
	if err != nil {
		return "", fmt.Errorf("marshal transaction RawData failed: %w", err)
	}
	txHash := sha256.Sum256(rawDataBytes)

	// 5️⃣ 核心修正：使用波场标准的 65 字节签名
	// 使用 SignCompact 生成 [V, R, S] 格式
	//signature := ecdsa.SignCompact(privKey, txHash[:], true)
	compactSig := ecdsa.SignCompact(privKey, txHash[:], true)
	v := (compactSig[0] - 27) & 3
	signature := make([]byte, 65)
	copy(signature[0:32], compactSig[1:33])   // 复制 R
	copy(signature[32:64], compactSig[33:65]) // 复制 S
	signature[64] = v                         // 将修正后的 V 放在最后
	// 6️⃣ 修正：将签名放入交易对象
	txExt.Transaction.Signature = [][]byte{signature}

	// 7️⃣ 广播交易
	result, err := grpcCli.Broadcast(txExt.Transaction)
	if err != nil {
		return "", fmt.Errorf("broadcast transaction failed: %w", err)
	}
	if !result.Result {
		// 打印具体错误原因，方便调试
		return "", fmt.Errorf("broadcast failed: %s", string(result.Message))
	}

	// 8️⃣ 获取交易 ID
	txid := hex.EncodeToString(txHash[:])

	return txid, nil
}

func (c *Client) TransferTRC20Back(
	from string,
	to string,
	contractAddr string,
	privateKeyHex string,
	amount *big.Int,
	feeLimit int64,
) (string, error) {
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return "error", err
	}
	fromAddr, _ := address.Base58ToAddress(from)
	toAddr, _ := address.Base58ToAddress(to)

	method := "transfer(address,uint256)"
	data := fmt.Sprintf(`[{"address": "%s"}, {"uint256": "0x%x"}]`, toAddr.String(), amount)
	// 2️⃣ 调用 TriggerContract 创建交易
	txExt, err := grpcCli.TriggerContract(
		fromAddr.String(), // 发起人地址，合约人地址
		contractAddr,      // 合约地址
		method,            // solidity 方法
		data,              // 参数 JSON
		feeLimit,          // feeLimit
		0,                 // call value, TRC20 通常是 0
		"",                // tTokenID
		0,                 // tTokenAmount
	)
	if err != nil {
		return "", fmt.Errorf("trigger TRC20 contract failed: %w", err)
	}

	// 3️⃣ 将私钥 hex 转 bytes
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	// 4️⃣ 使用 btcec 解析私钥
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	// 5️⃣ 获取交易 RawData bytes
	rawDataBytes, err := proto.Marshal(txExt.Transaction.RawData)
	if err != nil {
		return "", fmt.Errorf("交易人数据信息: %w", err)
	}

	// 6️⃣ 计算交易哈希
	txHash := sha256.Sum256(rawDataBytes)

	// 7️⃣ 使用私钥签名哈希（btcec/v2）
	sig := ecdsa.Sign(privKey, txHash[:]) // 返回 *ecdsa.Signature
	if sig == nil {
		return "", fmt.Errorf("使用签名信息")
	}
	// DER 编码签名
	signatureBytes := sig.Serialize()
	// 8️⃣ 添加签名到交易
	txExt.Transaction.Signature = append(
		txExt.Transaction.Signature,
		signatureBytes,
	)
	// 9️⃣ 广播交易
	result, err := grpcCli.Broadcast(txExt.Transaction)
	if err != nil {
		return "", fmt.Errorf("broadcast transaction failed: %w", err)
	}
	if !result.Result {
		return "", fmt.Errorf("broadcast TRC20 transaction failed")
	}

	// 🔟 获取交易 ID（txid）
	txidBytes := sha256.Sum256(rawDataBytes)
	txid := hex.EncodeToString(txidBytes[:])

	return txid, nil
}

// GetTransactionDetail 综合查询交易内容与执行状态
func (c *Client) GetTransactionDetail(txId string) (map[string]interface{}, error) {
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return nil, err
	}
	//
	log.Printf("[TRON-SCAN] 开始查询交易详情, TXID: %s", txId)
	tx, err := grpcCli.GetTransactionByID(txId)
	if err != nil || tx == nil {
		log.Printf("[TRON-SCAN] 链上未查到交易体或 RPC 报错, TXID: %s, Error: %v", txId, err)
		return nil, fmt.Errorf("没有查询到交易记录：%v", err)
	}
	info, err := grpcCli.GetTransactionInfoByID(txId)
	if err != nil {
		log.Printf("[TRON-SCAN] 获取 TransactionInfo 失败 (可能尚未打包), TXID: %s, Error: %v", txId, err)
	}
	detail := make(map[string]interface{})
	detail["txid"] = txId
	if len(tx.GetRawData().GetContract()) > 0 {
		contract := tx.GetRawData().GetContract()[0]
		detail["type"] = contract.GetType().String()
		detail["expiration"] = tx.GetRawData().GetExpiration()
		log.Printf("[TRON-SCAN] 交易类型: %s, 过期时间: %d", contract.GetType().String(), detail["expiration"])
	}
	//
	if info != nil {
		detail["status"] = info.GetResult().String()
		detail["block_height"] = info.BlockNumber
		detail["fee_usage"] = info.GetFee()
		log.Printf("[TRON-SCAN] 交易已入块: %d, 基础结果: %s, 消耗手续费: %d sun",
			info.BlockNumber, detail["status"], info.GetFee())
		if info.GetReceipt() != nil {
			receiptResult := info.GetReceipt().GetResult().String()
			detail["receipt_result"] = receiptResult
			detail["energy_usage"] = info.GetReceipt().GetEnergyUsage()
			detail["net_usage"] = info.GetReceipt().GetNetUsage()
			// 关键点：如果是 TRC20 转账失败，这里通常能看到原因
			if receiptResult != "DEFAULT" && receiptResult != "SUCCESS" {
				log.Printf("[TRON-SCAN] ⚠️ 警告：合约执行异常! ReceiptResult: %s, TXID: %s", receiptResult, txId)
			} else {
				log.Printf("[TRON-SCAN] 合约执行成功. 消耗能量: %d, 消耗带宽: %d",
					detail["energy_usage"], detail["net_usage"])
			}
		}
	} else {
		// 如果能查到 tx 但查不到 info，说明交易还在 Pending 中（未打包）
		detail["status"] = "PENDING"
		log.Printf("[TRON-SCAN] 交易状态：PENDING (已广播但尚未生成回执), TXID: %s", txId)
	}
	return detail, nil
}
