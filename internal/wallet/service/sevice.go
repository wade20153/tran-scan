package service

import (
	"context"
	"tron-scan/pkg/walletgrpc/walletpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletServiceServer struct {
	walletpb.UnimplementedWalletServiceServer
}

// GetUniqueId 实现
func (s *WalletServiceServer) GetUniqueId(ctx context.Context, req *walletpb.Empty) (*walletpb.UniqueIdResp, error) {
	return &walletpb.UniqueIdResp{UniqueId: 12345}, nil
}

// CreateUserWallet 实现
func (s *WalletServiceServer) CreateUserWallet(ctx context.Context, req *walletpb.CreateWalletReq) (*walletpb.CreateWalletResp, error) {
	wallet := &walletpb.WalletInfo{
		Chain:      "TRX",
		Address:    "T1234567890abcdef",
		PrivateKey: "xxxxxxxxxxxxxxxxxxxx",
		CreatedAt:  "2026-01-23T16:00:00Z",
	}
	return &walletpb.CreateWalletResp{
		Success: true,
		Wallet:  wallet,
	}, nil
}

// ListWallet 默认实现，返回未实现
func (s *WalletServiceServer) ListWallet(ctx context.Context, req *walletpb.MemberReq) (*walletpb.WalletListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWallet not implemented")
}

// GetETHWallet 默认实现，返回未实现
func (s *WalletServiceServer) GetETHWallet(ctx context.Context, req *walletpb.MemberReq) (*walletpb.WalletResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetETHWallet not implemented")
}

// GetTRXWallet 默认实现，返回未实现
func (s *WalletServiceServer) GetTRXWallet(ctx context.Context, req *walletpb.MemberReq) (*walletpb.WalletResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTRXWallet not implemented")
}

// TransferCrypto 默认实现，返回未实现
func (s *WalletServiceServer) TransferCrypto(ctx context.Context, req *walletpb.TransferReq) (*walletpb.TransferResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferCrypto not implemented")
}

// EstimateGas 默认实现，返回未实现
func (s *WalletServiceServer) EstimateGas(ctx context.Context, req *walletpb.EstimateGasReq) (*walletpb.EstimateGasResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimateGas not implemented")
}

// GetReconcileWallet 默认实现，返回未实现
func (s *WalletServiceServer) GetReconcileWallet(ctx context.Context, req *walletpb.Empty) (*walletpb.WalletResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReconcileWallet not implemented")
}

// TestTransferETH 默认实现，返回未实现
func (s *WalletServiceServer) TestTransferETH(ctx context.Context, req *walletpb.TransferReq) (*walletpb.TransferResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestTransferETH not implemented")
}

// TestTransferTRON 默认实现，返回未实现
func (s *WalletServiceServer) TestTransferTRON(ctx context.Context, req *walletpb.TransferReq) (*walletpb.TransferResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TestTransferTRON not implemented")
}

// GetBalance 默认实现，返回未实现
func (s *WalletServiceServer) GetBalance(ctx context.Context, req *walletpb.GetBalanceRequest) (*walletpb.GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
