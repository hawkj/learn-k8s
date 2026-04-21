# Kubernetes 自学教程（索引）

面向：本地 **kind** 或 **minikube** + 本仓库示例。文风简练；**每节末尾有练习题**，做完再进入下一节。

正文已按章拆分，请从下面入口顺序阅读。

## 章节目录

| 章 | 文件 |
|----|------|
| 第 0 章 环境与工具 | [chapters/chapter-00.md](./chapters/chapter-00.md) |
| 第 1 章 Gin 微服务 | [chapters/chapter-01.md](./chapters/chapter-01.md) |
| 第 2 章 容器化 | [chapters/chapter-02.md](./chapters/chapter-02.md) |
| 第 3 章 Kubernetes 核心对象 | [chapters/chapter-03.md](./chapters/chapter-03.md) |
| 第 4 章 可观测与排障 | [chapters/chapter-04.md](./chapters/chapter-04.md) |
| 第 5 章 CI/CD | [chapters/chapter-05.md](./chapters/chapter-05.md) |
| 第 6 章 多服务与调用链 | [chapters/chapter-06.md](./chapters/chapter-06.md) |
| 第 7 章 Service Mesh | [chapters/chapter-07.md](./chapters/chapter-07.md) |
| 第 8 章 弹性扩展 | [chapters/chapter-08.md](./chapters/chapter-08.md) |
| 第 9 章 加固与收尾 | [chapters/chapter-09.md](./chapters/chapter-09.md) |
| 附录 A / B | [chapters/appendix.md](./chapters/appendix.md) |

## 仓库结构（与正文对应）

| 路径 | 用途 |
|------|------|
| `learn-api/cmd/learn-api/` | 对外 API（健康检查、就绪、链式调用等） |
| `learn-api/cmd/backend/` | 第 6 章起：内部 backend |
| `learn-api/Dockerfile` / `Dockerfile.backend` | 第 2 章镜像 |
| `deploy/k8s/demo/` | 第 3 章起清单示例 |
| `.github/workflows/ci.yml` | 第 5 章 CI 示例 |
