package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// IPInfo IP 信息结构
type IPInfo struct {
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
	Country  string `json:"country"`
	Colo     string `json:"colo"`
	ASN      string `json:"asn"`
	Timezone string `json:"timezone"`
	Timestamp string `json:"timestamp"`
}

var version = "1.0.0"

func getIPInfo(r *http.Request) IPInfo {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	cfRay := getHeader(r, "CF-Ray", "")
	colo := ""
	asn := ""

	if cfRay != "" && len(cfRay) >= 9 {
		// CF-Ray 格式: xxxxxxxxxx-HKG (节点代码在连字符后)
		parts := strings.Split(cfRay, "-")
		if len(parts) == 2 {
			colo = parts[1]
		}
		// ASN 需要从其他头获取，这里暂时留空
		asn = ""
	}

	// 如果没有 CF 头，使用本地信息用于测试
	ipv4 := getHeader(r, "CF-Connecting-IP", "")
	if ipv4 == "" {
		ipv4 = getRemoteIP(r)
	}

	return IPInfo{
		IPv4:      ipv4,
		IPv6:      getHeader(r, "CF-Connecting-IPv6", ""),
		Country:   getHeader(r, "CF-IPCountry", ""),
		Colo:      colo,
		ASN:       asn,
		Timezone:  getHeader(r, "CF-Timezone", ""),
		Timestamp: now,
	}
}

func getRemoteIP(r *http.Request) string {
	ip := r.RemoteAddr
	// 去掉端口号
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	// 如果是 IPv6 映射的 IPv4，去掉前缀
	if strings.HasPrefix(ip, "[::1]") {
		return "127.0.0.1"
	}
	return ip
}

func getHeader(r *http.Request, key, defaultValue string) string {
	if val := r.Header.Get(key); val != "" {
		return val
	}
	return defaultValue
}

func htmlPage(info IPInfo, r *http.Request) string {
	host := r.Header.Get("Host")
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + host

	jsonBytes, _ := json.MarshalIndent(info, "", "  ")

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>IP Info Service</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
            padding: 40px 20px;
        }
        .container {
            max-width: 700px;
            margin: 0 auto;
        }
        .card {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            padding: 40px;
            margin-bottom: 20px;
        }
        h1 {
            color: #333;
            margin-bottom: 5px;
            font-size: 28px;
        }
        .version {
            color: #999;
            font-size: 14px;
            margin-bottom: 30px;
        }
        .section-title {
            color: #666;
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 15px;
            margin-top: 25px;
        }
        .section-title:first-of-type {
            margin-top: 0;
        }
        .info-item {
            padding: 15px;
            margin: 10px 0;
            background: #f8f9fa;
            border-radius: 8px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .label {
            color: #666;
            font-weight: 500;
        }
        .value {
            color: #333;
            font-weight: 600;
            font-family: "SF Mono", Monaco, monospace;
        }
        .value.empty {
            color: #ccc;
        }
        code {
            background: #1e1e1e;
            color: #d4d4d4;
            padding: 16px;
            border-radius: 8px;
            display: block;
            font-size: 13px;
            font-family: "SF Mono", Monaco, monospace;
            overflow-x: auto;
            margin: 10px 0;
        }
        .cmd {
            color: #4ec9b0;
        }
        .links {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
            margin-top: 10px;
        }
        .link {
            background: #f0f0f0;
            color: #333;
            text-decoration: none;
            padding: 8px 16px;
            border-radius: 6px;
            font-size: 14px;
            transition: background 0.2s;
        }
        .link:hover {
            background: #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <h1>IP Info Service</h1>
            <p class="version">Version %s</p>

            <p class="section-title">Your Information</p>
            <div class="info-item">
                <span class="label">IPv4</span>
                <span class="value">%s</span>
            </div>
            <div class="info-item">
                <span class="label">IPv6</span>
                <span class="value">%s</span>
            </div>
            <div class="info-item">
                <span class="label">Country</span>
                <span class="value">%s</span>
            </div>
            <div class="info-item">
                <span class="label">Data Center</span>
                <span class="value">%s</span>
            </div>
            <div class="info-item">
                <span class="label">Timezone</span>
                <span class="value">%s</span>
            </div>

            <p class="section-title">JSON API</p>
            <code><span class="cmd">curl -H "Accept: application/json" %s/</span></code>

            <p class="section-title">Response</p>
            <code>%s</code>

            <p class="section-title">Documentation</p>
            <div class="links">
                <a class="link" href="/llms.txt">llms.txt</a>
                <a class="link" href="/llm.txt">llm.txt</a>
                <a class="link" href="/robots.txt">robots.txt</a>
            </div>
        </div>
    </div>
</body>
</html>`, version, info.IPv4, displayOrEmpty(info.IPv6), displayOrEmpty(info.Country), displayOrEmpty(info.Colo), displayOrEmpty(info.Timezone), baseURL, string(jsonBytes))
}

func displayOrEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	info := getIPInfo(r)

	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-App-Version", version)
		json.NewEncoder(w).Encode(info)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-App-Version", version)
	fmt.Fprint(w, htmlPage(info, r))
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "User-agent: *\nAllow: /llms.txt\nAllow: /llm.txt\n")
}

func llmsTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	content := `# LLMs.txt - IP Info Service API Documentation

> IP Information Query Service - AI and LLM Friendly

This service provides visitor IP information via Cloudflare's request headers.

## API Endpoints

- GET / - Returns visitor IP information (JSON or HTML based on Accept header)
- GET /llms.txt - This file (LLMs.txt standard index)
- GET /llm.txt - Alternative LLM index (compatibility format)
- GET /robots.txt - Search engine rules

## Response Fields

| Field | Description |
|-------|-------------|
| ipv4 | Visitor's IPv4 address |
| ipv6 | Visitor's IPv6 address (if applicable) |
| country | ISO 3166-1 Alpha-2 country code |
| colo | Cloudflare edge node code |
| asn | Autonomous System Number |
| timezone | Visitor's timezone |
| timestamp | Request timestamp (ISO 8601) |

## Usage Examples

cURL:
  curl -H "Accept: application/json" https://your-worker.workers.dev/

JavaScript:
  const response = await fetch('https://your-worker.workers.dev/');
  const data = await response.json();
  console.log(data.country); // "CN"

Python:
  import requests
  response = requests.get('https://your-worker.workers.dev/')
  data = response.json()
  print(data['country'])  # "CN"

Sample Response:
  {
    "ipv4": "203.0.113.1",
    "ipv6": "",
    "country": "CN",
    "colo": "HKG",
    "asn": "13335",
    "timezone": "Asia/Shanghai",
    "timestamp": "2024-01-01T12:00:00.000Z"
  }

For more information, visit: https://github.com/xurenlu/ipinfo
`
	fmt.Fprint(w, content)
}

func llmTxtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	content := `# IP Info Service

A lightweight IP information query service based on Cloudflare Workers.

Quick Start:
  curl https://your-worker.workers.dev/

Links:
  Repository: https://github.com/xurenlu/ipinfo
  Full API Docs: /llms.txt
`
	fmt.Fprint(w, content)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", ipHandler)
	http.HandleFunc("/robots.txt", robotsHandler)
	http.HandleFunc("/llms.txt", llmsTxtHandler)
	http.HandleFunc("/llm.txt", llmTxtHandler)

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
