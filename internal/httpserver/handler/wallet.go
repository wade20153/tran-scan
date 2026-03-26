package handler

import (
	"encoding/json"
	"log"
	"net/http"
	request "tron-scan/internal/httpserver/handler/request"
	"tron-scan/internal/tronservice/service"
)

// WalletResponse 定义返回给前端的钱包结构体
type WalletResponse struct {
	Mnemonic    string `json:"mnemonic"`
	PrivateKey  string `json:"privateKey"`
	TronAddress string `json:"tronAddress"`
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址

	walletService := service.WalletService{}
	result, err := walletService.CreateUserWallet()
	if err != nil {
		log.Printf("创建钱包失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	resp := WalletResponse{
		Mnemonic:    result.Mnemonic,
		PrivateKey:  result.PrivateKey,
		TronAddress: result.TronAddress,
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("返回钱包信息失败: %v", err)
	}
}

func CreateWalletMain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址

	walletService := service.WalletService{}
	result, err := walletService.CreateUserMainWallet()
	if err != nil {
		log.Printf("创建钱包失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	resp := WalletResponse{
		Mnemonic:    result.Mnemonic,
		PrivateKey:  result.PrivateKey,
		TronAddress: result.TronAddress,
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("返回钱包信息失败: %v", err)
	}
}

// DeriveUserSubWallet 创建子钱包
func DeriveUserSubWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址
	var req request.DeriveSubWalletReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	walletService := service.WalletService{}
	result, err := walletService.DeriveUserSubWallet(req.Index)
	if err != nil {
		log.Printf("创建钱包失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	resp := WalletResponse{
		Mnemonic:    result.Mnemonic,
		PrivateKey:  result.PrivateKey,
		TronAddress: result.TronAddress,
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("返回钱包信息失败: %v", err)
	}
}

// QueryTronAccountAssets 创建子钱包
func QueryTronAccountAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址
	var req request.QueryWalletReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	walletService := service.WalletService{}
	result, err := walletService.QueryTronAccountAssets(req.Address)
	if err != nil {
		log.Printf("查询账户失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("返回钱包余额信息失败: %v", err)
	}
}

// GetTRC20TxHistory 获取所有请求
func GetTRC20TxHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址
	var req request.QueryWalletReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	walletService := service.WalletService{}
	result, err := walletService.GetTRC20TxHistory(req.Address)
	if err != nil {
		log.Printf("查询交易历史记录失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("返回钱包余额信息失败: %v", err)
	}
}

// TransferTRX 获取所有请求
func TransferTRX(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址
	var req request.TransferTRXReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "入参错误", http.StatusBadRequest)
		return
	}
	walletService := service.WalletService{}
	result, err := walletService.TransferTRX(req.FromAddr, req.ToAddr, req.Amount)
	if err != nil {
		log.Printf("转账TRX失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("转账失败: %v", err)
	}
}

// TransferTRC20 trc20转账
func TransferTRC20(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 调用创建波场链的地址
	var req request.TransferTRXReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "trc20转账入参错误", http.StatusBadRequest)
		return
	}
	walletService := service.WalletService{}
	result, err := walletService.TransferTRC20(req.FromAddr, req.ToAddr, req.Amount)
	if err != nil {
		log.Printf("转账trc20转账失败: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	// 设置返回头为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 写入 JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("trc20转账转账失败: %v", err)
	}
}
