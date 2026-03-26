package httpserver

import (
	"log"
	"net/http"
)

func Start(port string) *http.Server {
	if port == "" {
		log.Println("HTTP 服务未启用")
		return nil
	}

	mux := http.NewServeMux()
	RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP 服务启动成功，端口: %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 启动失败: %v", err)
		}
	}()

	return server
}
