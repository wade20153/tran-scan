package response

type TronAssetResult struct {
	Address       string // 钱包地址（T开头）
	TRX           int64  // TRX余额（单位：Sun，最小单位）
	TRXFloat      string // TRX余额（可读格式，保留6位小数）
	Bandwidth     int64  // 总带宽（免费+质押获得）用于普通转账
	BandwidthUsed int64  // 已使用带宽
	Energy        int64  // 总能量（质押获得）用于智能合约（TRC20转账）
	EnergyUsed    int64  // 已使用能量
	TRC20         map[string]string
}
