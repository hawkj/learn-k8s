# 第 5 章 CI/CD

**上一章**：[第 4 章](./chapter-04.md)｜**下一章**：[第 6 章 多服务](./chapter-06.md)

---

### 5.1 流水线阶段

**要点**：典型顺序：**lint** → **单元测试** → **构建二进制** → **构建镜像** → **推送仓库** → **更新集群**。

**本仓库**：`.github/workflows/ci.yml` 在 `learn-api` 下执行 `go test ./...`。

**练习 5.1**：将仓库推送到 GitHub，确认 Actions 中 **test** 任务通过（无 GitHub 时本地执行 `cd learn-api && go test ./...` 等价）。

---

### 5.2 镜像推送

**要点**：CI 中使用 `docker/login-action` 登录 **GHCR** / Docker Hub；镜像名建议 `ghcr.io/<owner>/learn-api:<git-sha>`。

**练习 5.2**：在流水线中增加 `docker build` + `push`（需仓库开启 Actions、配置 `GITHUB_TOKEN` 或 PAT），在制品库页面能看到新 tag。

---

### 5.3 部署到集群

**要点**：

- **方式 A**：`kubectl set image deployment/api api=<新镜像> -n demo`
- **方式 B**：Kustomize `images.newTag` 或 Helm `values`
- CI 需 **kubeconfig** 以 **Secret** 注入（注意安全）

**练习 5.3**：在 CI 或本机用可访问集群的 kubeconfig 执行一次 **apply** 或 **set image**，使 `api` 使用新镜像并成功滚动（本地 kind 可先手动 `kubectl set image` 验证）。

---

### 5.4 回滚

**要点**：`kubectl rollout history deployment/api -n demo`、`kubectl rollout undo deployment/api -n demo`。

**练习 5.4**：将 Deployment 镜像改为 `busybox` 等无法启动的应用，待 Pod 异常后执行 `rollout undo`，恢复上一版。
