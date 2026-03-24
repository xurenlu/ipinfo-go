# IP Info Service (Go 版本)

一个用 Go 语言实现的轻量级 IP 信息查询服务。

## 特性

- 纯 Go 实现，无需复杂的构建流程
- 零外部依赖（仅使用标准库）
- 支持 Cloudflare Workers 部署
- 提供 JSON 和 HTML 两种响应格式
- AI 友好（提供 llms.txt 和 llm.txt）

## 本地运行

```bash
cd go
go run main.go
```

默认运行在 `http://localhost:8080`

## 构建

```bash
go build -o ipinfo main.go
```

## 部署到 Cloudflare Workers

### 安装 Wrangler

```bash
npm install -g wrangler
```

### 创建 wrangler.toml

在 `go` 目录下创建 `wrangler.toml`:

```toml
name = "ip-info-go"
main = "worker.js"
compatibility_date = "2024-01-01"

[build]
command = "go build -o worker.js main.go"
```

### 部署

```bash
wrangler deploy
```

## API 使用

### 获取 JSON 响应

```bash
curl -H "Accept: application/json" http://localhost:8080/
```

### 响应示例

```json
{
  "ipv4": "203.0.113.1",
  "ipv6": "",
  "country": "CN",
  "colo": "HKG",
  "asn": "13335",
  "timezone": "Asia/Shanghai",
  "timestamp": "2024-01-01T12:00:00.000Z"
}
```

## 路由

| 路由 | 说明 |
|------|------|
| `GET /` | IP 信息（JSON/HTML） |
| `GET /robots.txt` | 搜索引擎规则 |
| `GET /llms.txt` | LLM 索引文件 |
| `GET /llm.txt` | LLM 索引文件（兼容） |
