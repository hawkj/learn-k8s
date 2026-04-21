# 第 2 章 容器化

**上一章**：[第 1 章](./chapter-01.md)｜**下一章**：[第 3 章 Kubernetes 核心](./chapter-03.md)

---

### 2.1 Dockerfile

**要点**：

- 多阶段：`golang` 构建 → `alpine` 运行，减小体积
- 非 root 用户（`appuser`）
- 本仓库：`learn-api/Dockerfile` 构建 **api**；`Dockerfile.backend` 构建 **backend**

**步骤（在项目根目录或 `learn-api` 下）：**

```bash
cd learn-api
docker build -t learn-api:local -f Dockerfile .
docker run --rm -p 8080:8080 learn-api:local
curl -s localhost:8080/healthz
```

**kind 加载本地镜像：**

```bash
kind load docker-image learn-api:local --name learn
```

**minikube：**

```bash
minikube image load learn-api:local
```

**练习 2.1**：容器内 `/healthz` 返回 `ok`；本地镜像已加载到集群节点（kind/minikube 按上表执行）。

---

### 2.2 .dockerignore

**要点**：缩小构建上下文、加快 `docker build`，避免把 `.git`、无关文件打进 context。

**练习 2.2**：对比 `docker build` 时 “Sending build context” 大小；在 `learn-api/.dockerignore` 中增加大目录忽略后再次构建，体积应减小或构建更快。
