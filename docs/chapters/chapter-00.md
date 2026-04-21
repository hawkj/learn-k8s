# 第 0 章 环境与工具

面向：本地 **kind** 或 **minikube** + 本仓库示例。每节末尾有练习题。

**上一章**：无｜**下一章**：[第 1 章 Gin 微服务](./chapter-01.md)

---

### 0.1 你需要什么

- **Go**：1.22+（与 `learn-api/go.mod` 一致）
- **Docker**：构建与运行镜像
- **kubectl**：操作集群
- **本地集群**：**kind** 或 **minikube** 二选一
- 可选：`helm`、`kustomize`、压测工具 **hey** 或 **wrk**

**要点**：本教程不绑定云厂商；后面 CI/CD、CA 等云相关小节可纸面完成。

**练习 0.1**：安装后执行 `kubectl version --client`、`docker version`，均能正常输出。

---

### 0.2 起一个本地集群

**kind 示例：**

```bash
kind create cluster --name learn
kubectl config use-context kind-learn   # 视 kind 输出为准
```

**minikube 示例：**

```bash
minikube start
```

**要点**：一个上下文对应一个集群；多集群时注意 `kubectl config current-context`。

**练习 0.2**：`kubectl get nodes` 中节点状态为 `Ready`。

---

### 0.3 命名约定

- **命名空间**：下文默认 `demo`（与 `deploy/k8s/demo/namespace.yaml` 一致）
- **镜像标签**：生产常用 **git commit 短 SHA**；本地可用 `local`

**练习 0.3**：执行 `kubectl apply -f deploy/k8s/demo/namespace.yaml`，或 `kubectl create namespace demo`，后续命令统一加 `-n demo`。
