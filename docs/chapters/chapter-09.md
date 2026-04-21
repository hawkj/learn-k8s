# 第 9 章 加固与收尾

**上一章**：[第 8 章](./chapter-08.md)｜**下一章**：[附录](./appendix.md)

---

### 9.1 NetworkPolicy

**要点**：默认允许全开；施加 **NetworkPolicy** 后仅声明的流量能进。**需 CNI 支持**（Calico、Cilium 等）；kind 默认 CNI 可能限制策略行为，以实际环境为准。

**步骤**：

```bash
kubectl apply -f deploy/k8s/demo/networkpolicy-backend.yaml
```

**练习 9.1**：从 **非 api** Pod（可临时 `kubectl run` 带 `app=other` 标签）访问 `backend` 应失败；从 **api** Pod 内访问应成功（若策略与 CNI 工作正常）。

---

### 9.2 检查清单

**要点**：镜像 tag 可追溯、探针齐全、requests/limits、Secret 不进 Git、回滚演练、HPA `maxReplicas` 有上限、生产开启 NetworkPolicy / RBAC。

**练习 9.2**：对照清单逐项勾选本仓库 demo，未完成项记入个人 TODO。
