package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tron-scan/config"
	"tron-scan/internal/httpserver"
	"tron-scan/internal/wallet/signer"
	"tron-scan/models"
	"tron-scan/pkg/db"

	walletGrpc "tron-scan/internal/walletgrpc"

	"google.golang.org/grpc"
)

func main() {

	log.Println("main 服务启动完成")
}

func init() {

	// 1.读取配置
	cfg, err := config.Load("local")
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	fmt.Printf("当前环境: %s, 日志等级: %s\n", cfg.App.Env, cfg.App.LogLevel)

	// 2.初始化 mysql 和 redis
	db.InitMySQL(cfg.MySQL)
	db.InitRedis(cfg.Redis)
	// 3 启动 启动 Wallet gRPC 服务
	grpcServer := walletGrpc.StartGrpcServer(config.GlobalConfig.GRPC)
	// 3. 启动 HTTP Server (可选, 例如提供 REST API)
	httpServer := startHTTPServer(config.GlobalConfig.App.HTTPPort)
	// 4. 加载
	// 4. 等待系统信号，优雅关闭
	waitForShutdown(grpcServer, httpServer)
	log.Println("环境初始化完成")

}

func loadSigner() {
	masterKey := []byte("1234567890123456") // AES 16字节示例
	store := signer.NewMemoryStore()
	manager := signer.NewManager(store, masterKey)
	// 假设管理员导入 xprv 已加密存储
	wallet := &models.Wallet{
		WalletID:      "wallet",
		Coin:          models.BTC,
		XPRVEncrypted: []byte{},
	}
	store.SaveWallet(wallet)
	// 启动时自动加载
	manager.LoadWallet("wallet")
	w := manager.GetWallet("")
	signer.InitAddressPool(w, 100)
	sig := signer.Sign(w, "helo")
	fmt.Printf(sig)
}

// startHTTPServer 启动 HTTP 服务 (可选)
func startHTTPServer(port string) *http.Server {
	if port == "" {
		log.Printf("HTTP 服务未启用")
		return nil
	}

	mux := http.NewServeMux()
	// 可在此注册 HTTP handler，例如 /health
	mux.HandleFunc("/health1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// 钱包相关
	httpserver.RegisterRoutes(mux)
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP 服务启动成功，监听端口: %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	return server
}

// waitForShutdown 等待系统信号并优雅关闭服务
func waitForShutdown(grpcServer *grpc.Server, httpServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Printf("接收到系统信号 %s, 开始优雅关闭服务...", sig)

	// 1. 关闭 gRPC
	if grpcServer != nil {
		log.Printf("停止 gRPC 服务...")
		grpcServer.GracefulStop()
	}

	// 2. 关闭 HTTP
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Printf("停止 HTTP 服务...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP 优雅关闭失败: %v", err)
		}
	}

	log.Printf("服务已优雅关闭")
}
