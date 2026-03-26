# 1.项目结构

```csharp
tron-wallet-service/
├── cmd/
│   └── server/
│       └── main.go                # 主启动文件
├── config/
│   ├── config.go                   # 配置加载
│   └── local.yml                   # 本地环境配置
├── internal/
│   ├── wallet/
│   │   ├── service/
│   │   │   └── wallet_service.go  # WalletServiceServer gRPC实现
│   │   └── signer/
│   │       └── signer.go           # 签名逻辑（TRX/TRC20）
│   └── grpc/
│       └── server.go               # gRPC 服务启动
├── pkg/
│   ├── rpc/
│   │   └── walletpb/               # protoc 生成的 gRPC 文件
│   └── tools/
│       └── crypto.go               # 助记词/私钥加解密工具
├── go.mod
└── go.sum
```



wsvip 的crypto-service，eth-crawler，tron-crawler   ，

 tron-crawler  的 crypto-service，eth-crawler，tron-crawler   ，

singer 
这个系统只对内，不对外，做签名

service 
这个系统虚拟币其他操作




一、整体架构总览（先把脑图立住）
```csharp
┌──────────────┐
│   phase.yml  │  ← 系统主密钥（不入库）
└──────┬───────┘
       │
       ▼
┌────────────────────┐
│   encryptor 模块   │  ← AES-256-GCM + PBKDF2
└──────┬─────────────┘
       │
       ▼
┌────────────────────┐
│   keystore 模块    │  ← 助记词 / key / salt 管理
└──────┬─────────────┘
       │
       ▼
┌────────────────────┐
│   wallet 模块      │  ← ETH / TRON 派生 + 归集
└────────────────────┘

pkg/
├── encryptor/     # 纯加解密（无业务）
├── keystore/      # 助记词 & key 管理
├── wallet/        # 派生地址 & 归集
├── config/        # phase.yml
└── security/      # 内存清理 / KMS / HSM
```

grpcs.Client：
只使用 gotron-sdk 官方 GrpcClient
✅ 不直接碰 proto / core / WalletClient（坑最多的地方）
✅ 自动管理 Start / Stop
✅ 支持 主网 / Nile 测试网
✅ 支持 TRX / TRC20 查询 & 转账
✅ 业务层完全无感
```sql
pkg/trx/grpcs/
├── client.go        # 客户端封装
├── balance.go       # 查询余额
├── transfer.go      # 转账
└── option.go        # 网络配置

```
你可以先只用 client.go + balance.go。