# 第 7 章 Service Mesh（选 Istio 或 Linkerd 之一）

**上一章**：[第 6 章](./chapter-06.md)｜**下一章**：[第 8 章 弹性扩展](./chapter-08.md)

---

### 7.1 安装控制面与注入

**要点**：按 [Istio](https://istio.io/latest/docs/setup/getting-started/) 或 [Linkerd](https://linkerd.io/2.14/getting-started/) 官方文档安装；对 `demo` 命名空间开启 **sidecar 自动注入**。

**练习 7.1**：`kubectl get pods -n demo -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{range .spec.containers[*]}{.name}{" "}{end}{"\n"}{end}'`，除业务容器外存在 **istio-proxy** 或 **linkerd-proxy**。

---

### 7.2 mTLS

**要点**：Mesh 内服务间流量默认 **双向 TLS**；严格模式需按文档开启 **PeerAuthentication** 等。

**练习 7.2**：使用所选 Mesh 的 dashboard 或 `istioctl proxy-config` / `linkerd viz` 等命令，确认服务间为加密流量。

---

### 7.3 流量管理（以 Istio 为例）

**要点**：`VirtualService` + `DestinationRule` 做 **权重分流**（如 backend `v1` 90% / `v2` 10%）；需为 **v2** 增加 Deployment 子集标签。

**练习 7.3**：部署 backend **v2**（复制 Deployment 改镜像/标签），配置 VS 权重，压测或统计请求分布。

---

### 7.4 与 CI/CD

**要点**：**发版/灰度** 往往改 **Git 清单或 Mesh 配置**；**HPA** 根据负载改 **副本数**，二者职责不同。

**练习 7.4**：用一两句话写清：本 demo 中「版本与流量比例」由谁管、「副本数」由谁管。
