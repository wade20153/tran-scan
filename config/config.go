package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// MySQLConfig 定义 MySQL 数据库的配置结构体
type MySQLConfig struct {
	Host            string `mapstructure:"host"`              // 数据库主机地址
	Port            int    `mapstructure:"port"`              // 数据库端口
	Username        string `mapstructure:"username"`          // 数据库用户名
	Password        string `mapstructure:"password"`          // 数据库密码（可加密后存储）
	Database        string `mapstructure:"database"`          // 数据库名称
	Charset         string `mapstructure:"charset"`           // 字符集，通常使用 utf8mb4
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`    // 数据库连接池最大空闲连接数
	MaxOpenConns    int    `mapstructure:"max_open_conns"`    // 数据库最大打开连接数
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // 连接最大存活时间（秒）
}

// RedisConfig 定义 Redis 的配置结构体
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`           // Redis 地址，例如 127.0.0.1:6379
	Password     string `mapstructure:"password"`       // Redis 密码
	DB           int    `mapstructure:"db"`             // Redis 数据库索引
	PoolSize     int    `mapstructure:"pool_size"`      // 连接池大小
	MinIdleConns int    `mapstructure:"min_idle_conns"` // 最小空闲连接数
}

// GRPCServerConfig 定义 gRPC 服务配置
type GRPCServerConfig struct {
	Port      string `mapstructure:"port"`       // gRPC 服务监听端口
	EnableTLS bool   `mapstructure:"enable_tls"` // 是否启用 TLS
	TLSCert   string `mapstructure:"tls_cert"`   // TLS 证书路径
	TLSKey    string `mapstructure:"tls_key"`    // TLS 私钥路径
}

// AppConfig 定义应用整体配置结构体
type AppConfig struct {
	App        App              `mapstructure:"app"`   // 解析 app 分组
	MySQL      MySQLConfig      `mapstructure:"mysql"` // MySQL 配置
	Redis      RedisConfig      `mapstructure:"redis"` // Redis 配置
	GRPC       GRPCServerConfig `mapstructure:"grpcs"` // gRPC 服务配置
	Tron       TronConfig       `mapstructure:"tron"`
	BlockChain BlockChain       `toml:"block_chain"` //区块链
}
type App struct {
	Env      string `mapstructure:"env"`
	LogLevel string `mapstructure:"log_level"`
	HTTPPort string `mapstructure:"http_port"` // 对应 app.http_port
}

// TronConfig 定义 TRON 链相关配置
type TronConfig struct {
	Network      string   `mapstructure:"network"`       // 网络类型: nile / mainnet
	FullNode     string   `mapstructure:"fullnode"`      // FullNode RPC
	SolidityNode string   `mapstructure:"soliditynode"`  // SolidityNode RPC
	EventServer  string   `mapstructure:"eventserver"`   // 事件服务器
	USDTContract string   `mapstructure:"usdt_contract"` // USDT(TRC20) 合约地址
	TronApiUrl   string   `mapstructure:"tron_url"`      //tron api key
	TronHttp     string   `mapstructure:"tron_http"`     //tron api key
	TronApiKey   []string `mapstructure:"tron_api_key"`  //tron api key
}

type BlockChain struct {
	SignApiUrl          string   `toml:"sign_api_url"`          //签名服务
	SignPubKey          string   `toml:"sign_pub_key"`          //签名公钥
	SignPlatCode        string   `toml:"sign_plat_code"`        //平台编号
	TronApiUrl          []string `toml:"tron_api_url"`          //tron grpc 地址
	TronApiKey          []string `toml:"tron_api_key"`          //tron api key
	TronContractAddress string   `toml:"tron_contract_address"` //tron 合约地址
	TronMinGasAmount    string   `toml:"tron_min_gas_amount"`   //归集需要的最小手续费
	TronMinUsdtAmount   string   `toml:"tron_min_usdt_amount"`  //自动归集最小余额
	EthApiUrl           []string `toml:"eth_api_url"`           //eth grpc 地址
	EthContractAddress  string   `toml:"eth_contract_address"`  //eth 合约地址
	EthMinGasAmount     string   `toml:"eth_min_gas_amount"`    //归集需要的最小手续费
	EthMinUsdtAmount    string   `toml:"eth_min_usdt_amount"`   //自动归集最小余额
}

// GlobalConfig 用于存储全局配置信息
var GlobalConfig AppConfig

// Load 从指定环境配置文件加载配置
// env: 配置文件名，不带扩展名，例如 dev / prod
// 返回值: AppConfig 配置结构体和 error
func Load(env string) (AppConfig, error) {
	// 设置配置文件名称和类型
	viper.SetConfigName(env)
	viper.SetConfigType("yml")

	// 配置文件搜索路径（按顺序查找）
	viper.AddConfigPath("./config")        // 当前目录 ./config
	viper.AddConfigPath("../../config")    // 项目上级目录
	viper.AddConfigPath("../../../config") // 更上级目录，可适配不同项目结构

	conf := AppConfig{}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将配置解析到局部结构体
	if err := viper.Unmarshal(&conf); err != nil {
		return conf, fmt.Errorf("解析配置到结构体失败: %w", err)
	}

	// 将配置解析到全局结构体，方便全局调用
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return conf, fmt.Errorf("解析配置到全局结构体失败: %w", err)
	}

	return conf, nil
}
