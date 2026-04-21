# 第 8 章 弹性扩展

**上一章**：[第 7 章](./chapter-07.md)｜**下一章**：[第 9 章 加固与收尾](./chapter-09.md)

---

### 8.1 Metrics Server

**要点**：`kubectl top`、`HPA` 依赖 **metrics-server**。kind 常需安装：

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# 若 TLS 问题，可对 deployment 增加 --kubelet-insecure-tls（仅本地）
```

**练习 8.1**：`kubectl top nodes` 与 `kubectl top pods -n demo` 有数值（需 Pod 运行一段时间）。

---

### 8.2 HPA

**要点**：`HorizontalPodAutoscaler` 关联 **Deployment**，依据 CPU/内存或自定义指标扩缩。**CPU 利用率** 需容器设了 `resources.requests.cpu`。

**步骤**：`kubectl apply -f deploy/k8s/demo/hpa-api.yaml`，对 `api` Service 加压：

```bash
# 宿主机 port-forward 后
hey -n 50000 -c 50 http://127.0.0.1:8080/api/v1/hello
```

**练习 8.2**：观察 `kubectl get hpa -n demo` 的 **TARGETS** 与 **REPLICAS** 变化；降压后副本逐步回落。

---

### 8.3 Cluster Autoscaler（云环境）

**要点**：节点不足时 **CA** 扩容节点组；节点闲置时缩容。**本地 kind 一般不部署 CA**。

**练习 8.3**（纸面）：画出 **Pod Pending** → CA 增加节点 → 调度成功；利用率低 → CA 缩节点的流程。

---

### 8.4 VPA（进阶）

**要点**：**VPA** 调 **requests/limits**；与 **HPA** 同时用于同一 workload 时需读[官方说明](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler)，避免打架。

**练习 8.4**（可选）：写出三条：何时优先 **HPA**、何时考虑 **VPA**、何时只用其一。
