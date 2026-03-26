package grpcs

import (
	"encoding/hex"
	_ "encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
	"tron-scan/config"
	"tron-scan/internal/httpserver/handler/response"

	_ "github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	_ "github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	
	_ "github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
)

// GetTrxBalance 查询 TRX 余额（返回 SUN）
func (c *Client) GetTrxBalance(addr string) (*big.Int, error) {
	// 3️⃣ 获取 client（关键）
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return nil, err
	}
	// 3️⃣ 调用 TriggerConstantContract
	acc, err := grpcCli.GetAccount(addr)
	if err != nil {
		return nil, err
	}
	return big.NewInt(acc.Balance), nil
}

// GetTrc20Balance 查询 TRC20 余额（返回最小单位，如 USDT = 6 位）
func (c *Client) GetTrc20Balance(userAddr string, contractAddr string) (*big.Int, error) {
	// 1️⃣ Base58 地址检查
	fromAddr, err := address.Base58ToAddress(userAddr)
	if err != nil {
		fmt.Printf("无效的用户地址: %s, err=%v\n", userAddr, err)
		return nil, fmt.Errorf("无效的用户地址: %w", err)
	}

	_, err = address.Base58ToAddress(contractAddr)
	if err != nil {
		fmt.Printf("无效的合约地址: %s, err=%v\n", contractAddr, err)
		return nil, fmt.Errorf("无效的合约地址: %w", err)
	}

	// 2️⃣ JSON 参数（address value 必须是 bytes 对象）
	paramJSON := fmt.Sprintf(`[{"type":"address","value":{"bytes":"%s"}}]`, hex.EncodeToString(fromAddr.Bytes()))
	data := fmt.Sprintf(`[{"address": "%s"}]`, fromAddr.String())
	fmt.Printf("TriggerConstantContract 参数: user=%s, contract=%s, method=balanceOf(address), params=%s\n",
		userAddr, contractAddr, paramJSON)
	// 3️⃣ 获取 client（关键）
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return nil, err
	}
	// 3️⃣ 调用 TriggerConstantContract

	txExt, err := grpcCli.TriggerConstantContract(
		"",
		contractAddr,
		"balanceOf(address)",
		data,
	)
	if err != nil {
		fmt.Printf("TriggerConstantContract 调用失败, from=%s, contract=%s, method=balanceOf(address), err=%v\n",
			userAddr, contractAddr, err)
		c.ReportFail(grpcCli)
		return nil, fmt.Errorf("TriggerConstantContract 调用失败: %w", err)
	}
	// ✅ 成功上报
	c.ReportSuccess(grpcCli)
	// 4️⃣ 解析返回值
	if len(txExt.ConstantResult) == 0 {
		fmt.Printf("TriggerConstantContract 返回结果为空\n")
		return big.NewInt(0), nil
	} else {
		// 打印第一个 ConstantResult
		fmt.Printf("ConstantResult[0] (hex): %x\n", txExt.ConstantResult[0])
	}
	// 打印整个 txExt
	txExtJSON, _ := json.MarshalIndent(txExt, "", "  ")
	fmt.Printf("完整 txExt:\n%s\n", string(txExtJSON))
	// 解析
	balance := new(big.Int).SetBytes(txExt.ConstantResult[0])
	fmt.Printf("查询到 TRC20 余额: %s\n", balance.String())
	return balance, nil
}

// GetAccountResource  查询资源
func (c *Client) GetAccountResource(addr string) (*api.AccountResourceMessage, error) {
	// 3️⃣ 获取 client（关键）
	grpcCli, err := c.GetOneClient()
	if err != nil {
		return nil, err
	}
	// 3️⃣ 调用 TriggerConstantContract
	resource, err := grpcCli.GetAccountResource(addr)
	if err != nil {
		return nil, err
	}

	return resource, nil

}

func (c *Client) GetTRC20TxHistory(userAddr string, fingerprint string) (*response.TRC20HistoryResult, error) {
	baseURL := config.GlobalConfig.Tron.TronHttp
	if baseURL == "" {
		baseURL = "https://nile.trongrid.io"
	}

	// 1. 组装 URL，带上 fingerprint 参数
	apiPath := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20?limit=20&only_confirmed=true",
		strings.TrimSuffix(baseURL, "/"), userAddr)

	if fingerprint != "" {
		apiPath = fmt.Sprintf("%s&fingerprint=%s", apiPath, fingerprint)
	}

	// 2. 发起请求 (逻辑同前，省略部分重复代码)
	req, _ := http.NewRequest("GET", apiPath, nil)
	// 设置 API Key...

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 3. 解析
	var tronResp response.TronGridResponse
	if err := json.NewDecoder(resp.Body).Decode(&tronResp); err != nil {
		return nil, err
	}

	// 4. 数据后处理
	for _, tx := range tronResp.Data {
		if tx == nil {
			continue
		}
		val := new(big.Int)
		val.SetString(tx.RawValue, 10)
		tx.Value = val
	}

	// 5. 封装返回结果
	result := &response.TRC20HistoryResult{
		Transfers:   tronResp.Data,
		Fingerprint: tronResp.Meta.Fingerprint,
		HasMore:     tronResp.Meta.Fingerprint != "", // 有指纹代表还有下一页
		PageSize:    len(tronResp.Data),
	}

	return result, nil
}
