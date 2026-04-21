# 第 1 章 Gin 微服务（应用本体）

**上一章**：[第 0 章](./chapter-00.md)｜**下一章**：[第 2 章 容器化](./chapter-02.md)

---

### 1.1 最小 HTTP 服务

**要点**：

- `GET /healthz`：进程存活即 200，正文 `ok`
- `GET /api/v1/hello?name=xx`：返回 JSON 问候（见 `learn-api/cmd/learn-api/main.go`）

**步骤：**

```bash
cd learn-api
go run ./cmd/learn-api
# 另开终端
curl -s localhost:8080/healthz
curl -s 'localhost:8080/api/v1/hello?name=k8s'
```

**练习 1.1**：本机 `curl` 上述两路径均符合预期。

---

### 1.2 存活 vs 就绪

**要点**：

- **Liveness**：不通过则 **重启容器**（卡住、死锁时有用）
- **Readiness**：不通过则 **从 Service 摘流**（依赖未就绪时有用）
- 本示例：`/readyz` 在进程启动 **约 3 秒** 后才返回 200，此前为 503（模拟依赖加载）

**练习 1.2**：启动后立即 `curl -i localhost:8080/readyz` 为 503；3 秒后为 200。

---

### 1.3 配置

**要点**：端口、日志级别由环境变量控制：`PORT`（默认 `8080`）、`LOG_LEVEL`（`debug`/`info`/…）。

**练习 1.3**：`PORT=9090 go run ./cmd/learn-api`，对 `9090` 端口访问 `/healthz`；再换端口启动第二实例（需不同 `PORT`），两进程互不冲突。

---

### 1.4 日志

**要点**：使用 **stdout** 上的 **JSON 日志**（`log/slog`），便于采集与检索。

**练习 1.4**：访问 `/api/v1/hello`，在运行 `go run` 的终端看到带 `path`、`name` 的 JSON 日志行。

---

### 1.5 优雅退出

**要点**：收到 `SIGTERM`/`SIGINT` 后，在 **10 秒** 内调用 `http.Server.Shutdown`，结束在途请求。

**步骤**：先请求 `/slow`（约 8 秒），在请求进行中对该进程 `kill -15 <pid>`，应看到请求仍能完成或按超时结束，而非直接断连。

**练习 1.5**：用 `curl` 访问 `/slow` 的同时发 `SIGTERM`，确认进程在关闭前尽量处理完（或观察日志行为符合预期）。
