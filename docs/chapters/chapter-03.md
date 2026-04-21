# 第 3 章 Kubernetes 核心对象

**上一章**：[第 2 章](./chapter-02.md)｜**下一章**：[第 4 章 可观测与排障](./chapter-04.md)

---

### 3.1 Deployment

**要点**：`Deployment` 管理 Pod 副本与滚动更新；镜像需集群 **可拉取**（本地集群用 `kind load` / `minikube image load`）。

**步骤：**

```bash
kubectl apply -f deploy/k8s/demo/namespace.yaml
kubectl apply -f deploy/k8s/demo/configmap-app.yaml
kubectl apply -f deploy/k8s/demo/secret-app.yaml
kubectl apply -f deploy/k8s/demo/deployment-api.yaml
kubectl get pods -n demo -l app=api
```

**练习 3.1**：`demo` 命名空间下 **2** 个 `api` Pod 均为 `Running`。

---

### 3.2 Service（ClusterIP）

**要点**：`Service` 通过 **selector** 匹配 Pod，提供稳定 **DNS**：`api.demo.svc.cluster.local`（同命名空间内可简写 `api`）。

**步骤：**

```bash
kubectl apply -f deploy/k8s/demo/service-api.yaml
kubectl run curl --rm -it --restart=Never -n demo --image=curlimages/curl -- \
  curl -s http://api:8080/healthz
```

**练习 3.2**：集群内 curl 访问 `http://api:8080/healthz` 返回 `ok`。

---

### 3.3 ConfigMap / Secret

**要点**：配置与镜像分离；敏感数据进 **Secret**，勿写进镜像或 Git（示例 `secret-app.yaml` 仅作本地演示）。

**练习 3.3**：修改 `deploy/k8s/demo/configmap-app.yaml` 中 `LOG_LEVEL` 为 `debug`，`kubectl apply` 后滚动重启 Pod（或 `kubectl rollout restart deployment/api -n demo`），日志级别应变化（需结合 `slog` 行为观察）。

---

### 3.4 探针

**要点**：清单中 `livenessProbe` → `/healthz`，`readinessProbe` → `/readyz`。

**练习 3.4**：`kubectl describe pod -n demo -l app=api`，查看 **Ready** 条件；启动后前几秒 readiness 可能失败，就绪后 **Endpoints** 才包含该 Pod（`kubectl get endpoints -n demo api`）。

---

### 3.5 Ingress（可选）

**要点**：HTTP 入口需集群内 **Ingress Controller**（kind 常配合 [ingress-nginx](https://kind.sigs.k8s.io/docs/user/ingress/)）；minikube 可用 `minikube addons enable ingress`。

**练习 3.5**（可选）：按所选环境文档安装 Controller，创建 `Ingress` 资源，从宿主机用域名访问到 `api` Service。

---

### 3.6 资源 requests/limits

**要点**：`requests` 供调度与 HPA；`limits` 防止单 Pod 占满节点。过小会 **OOMKilled** 或 CPU **节流**。

**练习 3.6**：将 `deployment-api.yaml` 中 `memory.limits` 调到极低（如 `32Mi`），重新 apply，观察 Pod 是否 **OOM**（`kubectl describe pod` 中 `Last State`）。
