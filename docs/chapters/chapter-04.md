# 第 4 章 可观测与排障

**上一章**：[第 3 章](./chapter-03.md)｜**下一章**：[第 5 章 CI/CD](./chapter-05.md)

---

### 4.1 常用命令

**要点**：

| 命令 | 用途 |
|------|------|
| `kubectl logs [-f] pod/...` | 容器日志 |
| `kubectl describe pod/...` | 事件与状态 |
| `kubectl get events -n demo --sort-by=.lastTimestamp` | 排障时间线 |
| `kubectl port-forward -n demo svc/api 8080:8080` | 本地访问 Service |
| `kubectl exec -it pod/... -- sh` | 进容器 |

**练习 4.1**：`kubectl port-forward` 访问 `api`，再故意将镜像改为无效并滚动发布，用 `describe` + `events` 定位 **ImagePullBackOff** 或 **CrashLoop**。

---

### 4.2 指标（可选）

**要点**：Prometheus 拉取 Pod 指标需 **Pod 暴露 metrics** 且 **ServiceMonitor** 等（视安装方式而定）；本章可只读概念。

**练习 4.2**（可选）：安装 [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts) 或简易 Prometheus，能查询到集群基础指标即可。
