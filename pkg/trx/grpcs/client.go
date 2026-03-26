package grpcs

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"tron-scan/config"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
)

var (
	globalClient *Client
	onceClient   sync.Once
)

type grpcWrapper struct {
	cli       *client.GrpcClient
	failCount int32
	available int32 // 1=可用 0=不可用
	apiKey    string
}

type Client struct {
	clients []*grpcWrapper
	index   uint32
	mu      sync.RWMutex
}

// ================== 单例 ==================

func GetClient() (*Client, error) {
	var err error

	onceClient.Do(func() {
		globalClient = &Client{}
		err = globalClient.init()
		go globalClient.healthCheckLoop()
	})

	return globalClient, err
}

// ================== 初始化 ==================

func (c *Client) init() error {
	apiUrl := config.GlobalConfig.Tron.TronApiUrl
	apiKeys := config.GlobalConfig.Tron.TronApiKey

	if len(apiKeys) == 0 {
		return fmt.Errorf("no api keys configured")
	}

	for _, apiKey := range apiKeys {
		cli := client.NewGrpcClient(apiUrl)
		cli.SetTimeout(time.Second * 60)

		if apiKey != "" {
			cli.SetAPIKey(apiKey)
		}

		if err := cli.Start(client.GRPCInsecure()); err != nil {
			continue
		}

		wrapper := &grpcWrapper{
			cli:       cli,
			failCount: 0,
			available: 1,
			apiKey:    apiKey,
		}

		c.clients = append(c.clients, wrapper)
	}

	if len(c.clients) == 0 {
		return fmt.Errorf("no available grpc clients")
	}

	return nil
}

// ================== 轮询获取 ==================

func (c *Client) GetOneClient() (*client.GrpcClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	n := len(c.clients)
	if n == 0 {
		return nil, fmt.Errorf("no clients")
	}

	// 轮询 index
	start := atomic.AddUint32(&c.index, 1)

	for i := 0; i < n; i++ {
		idx := int((start + uint32(i)) % uint32(n))
		w := c.clients[idx]

		if atomic.LoadInt32(&w.available) == 1 {
			return w.cli, nil
		}
	}

	return nil, fmt.Errorf("no available grpc client")
}

// ================== 上报失败 ==================

func (c *Client) ReportFail(cli *client.GrpcClient) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, w := range c.clients {
		if w.cli == cli {
			fail := atomic.AddInt32(&w.failCount, 1)

			// 连续失败3次 → 标记不可用
			if fail >= 3 {
				atomic.StoreInt32(&w.available, 0)
			}
			return
		}
	}
}

// ================== 上报成功 ==================

func (c *Client) ReportSuccess(cli *client.GrpcClient) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, w := range c.clients {
		if w.cli == cli {
			atomic.StoreInt32(&w.failCount, 0)
			atomic.StoreInt32(&w.available, 1)
			return
		}
	}
}

// ================== 健康检查 ==================

func (c *Client) healthCheckLoop() {
	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		c.mu.RLock()
		for _, w := range c.clients {
			if atomic.LoadInt32(&w.available) == 1 {
				continue
			}

			// 尝试恢复
			err := w.cli.Start(client.GRPCInsecure())
			if err == nil {
				atomic.StoreInt32(&w.available, 1)
				atomic.StoreInt32(&w.failCount, 0)
			}
		}
		c.mu.RUnlock()
	}
}

// ================== 关闭 ==================

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, w := range c.clients {
		if w.cli != nil {
			w.cli.Stop()
		}
	}
}
