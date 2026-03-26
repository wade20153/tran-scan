package response

import "math/big"

type TRC20Transfer struct {
	TransactionID string   `json:"transaction_id"`
	From          string   `json:"from"`
	To            string   `json:"to"`
	Value         *big.Int `json:"-"`     // 用于内部逻辑
	RawValue      string   `json:"value"` // 接口返回的是字符串形式的数字
	Type          string   `json:"type"`
	TokenInfo     struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
	BlockTimestamp int64 `json:"block_timestamp"`
}

type TronGridResponse struct {
	Data    []*TRC20Transfer `json:"data"`
	Success bool             `json:"success"`
	Meta    struct {
		Fingerprint string `json:"fingerprint"`
	} `json:"meta"`
}

// 业务层通用的分页结果包装
type TRC20HistoryResult struct {
	Transfers   []*TRC20Transfer `json:"transfers"`
	Fingerprint string           `json:"fingerprint"` // 下一页的凭证
	HasMore     bool             `json:"has_more"`    // 是否还有更多
	PageSize    int              `json:"page_size"`   // 本次返回的数量
}
