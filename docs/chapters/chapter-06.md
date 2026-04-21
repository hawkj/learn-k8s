# 第 6 章 多服务与调用链

**上一章**：[第 5 章](./chapter-05.md)｜**下一章**：[第 7 章 Service Mesh](./chapter-07.md)

---

### 6.1 拆服务

**要点**：**api** 对外；**backend** 仅集群内访问。api 通过环境变量 `BACKEND_BASE_URL` 调用 backend（见 `main.go` 中 `/api/v1/chain`）。

**步骤：**

```bash
cd learn-api
docker build -t learn-backend:local -f Dockerfile.backend .
kind load docker-image learn-backend:local --name learn   # 或 minikube image load

kubectl apply -f deploy/k8s/demo/deployment-backend.yaml
kubectl apply -f deploy/k8s/demo/service-backend.yaml
kubectl set env deployment/api -n demo BACKEND_BASE_URL=http://backend.demo.svc.cluster.local:8080
```

等待 `api` Pod 重启后：

```bash
kubectl run curl --rm -it --restart=Never -n demo --image=curlimages/curl -- \
  curl -s http://api:8080/api/v1/chain
```

**练习 6.1**：响应 JSON 中含 `via":"api"` 且 `backend_status` 为 200。

---

### 6.2 统一标签与选择器

**要点**：`app`、`version` 标签便于 **Service / Mesh / NetworkPolicy** 选择同一组 Pod。

**练习 6.2**：`kubectl get pods -n demo --show-labels`，确认 **backend** Pod 带 `version=v1`。
