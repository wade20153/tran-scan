package httpserver

import (
	"net/http"
	"tron-scan/internal/httpserver/handler"
)

func RegisterRoutes(mux *http.ServeMux) {

	// 基础
	mux.HandleFunc("/health", handler.Health)

	// 钱包
	mux.HandleFunc("/wallet/create", handler.CreateWallet)
	// 钱包
	mux.HandleFunc("/wallet/create/mian", handler.CreateWalletMain)

	// 钱包
	mux.HandleFunc("/wallet/create/sub", handler.DeriveUserSubWallet)
	// QueryTronAccountAssets
	// 查询账户 TRX + TRC20 资产，其他资产
	mux.HandleFunc("/wallet/banlance", handler.QueryTronAccountAssets)
	//mux.HandleFunc("/wallet/faucet", handler.RequestTestnetTRX)      // Nile 测试网领取 TRX
	//mux.HandleFunc("/wallet/faucet/trc20", handler.RequestTestnetTRC20)// Nile 测试网领取 TRC20

	// -----------------------------
	// 账户管理
	// -----------------------------
	//mux.HandleFunc("/account/balance", handler.GetBalance)           // 查询账户余额
	//mux.HandleFunc("/account/assets", handler.GetAccountAssets)      // 查询账户所有资产
	mux.HandleFunc("/account/history", handler.GetTRC20TxHistory) // 查询账户交易历史

	// -----------------------------
	// 链操作 / 交易
	// -----------------------------
	mux.HandleFunc("/chain/transfer/trx", handler.TransferTRX)     // 转账 TRX
	mux.HandleFunc("/chain/transfer/trc20", handler.TransferTRC20) // 转账 TRC20 Token
	//mux.HandleFunc("/chain/transaction/status", handler.GetTransactionStatus)// 查询交易状态
	//mux.HandleFunc("/chain/contract/call", handler.CallContract)     // 调用智能合约
	//mux.HandleFunc("/chain/contract/deploy", handler.DeployContract) // 部署智能合约
	//mux.HandleFunc("/wallet/estimateFee", handler.EstimateFee)       // 估算交易手续费
	//mux.HandleFunc("/wallet/energy/bandwidth", handler.GetEnergyAndBandwidth) // 查询能量/带宽
	//mux.HandleFunc("/wallet/transaction/list", handler.ListTransactions)      // 查询交易列表
	//mux.HandleFunc("/wallet/transaction/details", handler.GetTransactionDetails) // 查询交易详情
}
