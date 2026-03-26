package walletgrpc

import (
	"log"
	"net"
	"tron-scan/config"
	"tron-scan/internal/wallet/service"
	"tron-scan/pkg/walletgrpc/walletpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ServerConfig gRPC 服务配置
type ServerConfig struct {
	Port string
}

// StartWalletServer 启动 Wallet gRPC 服务
func StartGrpcServer(grpcCfg config.GRPCServerConfig) *grpc.Server {
	var opts []grpc.ServerOption

	// TLS 配置
	if grpcCfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(grpcCfg.TLSCert, grpcCfg.TLSKey)
		if err != nil {
			log.Fatalf("加载 TLS 证书失败: %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	// 注册 WalletServiceServer 也可以注册其他服务
	walletpb.RegisterWalletServiceServer(grpcServer, &service.WalletServiceServer{})
	log.Printf("WalletService gRPC 已注册")
	// 监听 TCP 端口
	listener, err := net.Listen("tcp", grpcCfg.Port)
	if err != nil {
		log.Fatalf("gRPC 监听失败: %v", err)
	}

	log.Printf("Wallet gRPC server is running on %s...", grpcCfg.Port)
	go func() {
		log.Printf("gRPC 服务启动成功，监听端口: %s", grpcCfg.Port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC 服务启动失败: %v", err)
		}
	}()
	return grpcServer
}
