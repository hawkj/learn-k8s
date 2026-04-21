# 第 2 章 容器化

**上一章**：[第 1 章](./chapter-01.md)｜**下一章**：[第 3 章 Kubernetes 核心](./chapter-03.md)

---

### 2.1 Dockerfile

**要点**：

- 多阶段：`golang` 构建 → `alpine` 运行，减小体积
- 非 root 用户（`appuser`）
- 本仓库：`learn-api/Dockerfile` 构建 **api**；`Dockerfile.backend` 构建 **backend`

**Dockerfile 逐行（`learn-api/Dockerfile`）**

镜像分 **两阶段**：前一阶段用 Go 工具链**编译**出二进制；后一阶段只用 **Alpine** 装运行所需文件，最终镜像里**不带** Go SDK 和源码，体积更小。

| 行号（文件） | 指令 | 作用 |
|----|------|------|
| 1 | `# syntax=docker/dockerfile:1` | 声明 Dockerfile 语法版本，便于 BuildKit 使用较新的解析能力（可选但推荐）。 |
| 2 | `FROM golang:1.22-alpine AS build` | **第一阶段**命名为 `build`：基于带 Go 1.22 的 Alpine，后续可用 `COPY --from=build` 引用该阶段产物。 |
| 3 | `WORKDIR /src` | 之后 `RUN`/`COPY` 的默认工作目录设为 `/src`。 |
| 4 | `RUN apk add --no-cache ca-certificates git` | 安装 `go mod` 可能需要的依赖（HTTPS、部分模块拉取需要 `git`）；`--no-cache` 减小层体积。 |
| 5 | `COPY go.mod go.sum ./` | **先只复制依赖清单**，与下面 `go mod download` 组成一层，源码变更时仍可复用依赖层缓存。 |
| 6 | `RUN go mod download` | 下载模块，不整盘 `COPY` 前先拉依赖，**加速重复构建**。 |
| 7 | `COPY . .` | 再复制其余源码与资源。 |
| 8 | `RUN CGO_ENABLED=0 go build ...` | 编译静态二进制（见下「关键：`go build` 一行」）。 |
| 9 | （空行） | 分隔两阶段，无指令。 |
| 10 | `FROM alpine:3.20` | **第二阶段**：全新基础镜像，仅作运行环境，不含上一阶段的编译工具。 |
| 11–12 | `RUN apk add ... && adduser ...` | 安装 CA 证书；创建 **非 root** 用户 `appuser`（uid 65532，与常见 restricted 约定一致）。 |
| 13 | `WORKDIR /app` | 进程工作目录。 |
| 14 | `COPY --from=build /out/learn-api /app/learn-api` | 仅从 `build` 阶段拷贝**编译产物**，不把源码和 `GOMODCACHE` 打进最终镜像。 |
| 15 | `USER appuser` | 之后进程以非特权用户运行，降低容器内被利用后的风险。 |
| 16 | `EXPOSE 8080` | **声明**容器对外端口（文档/约定，**不会**自动映射主机端口；`docker run -p` 才映射）。 |
| 17 | `ENV PORT=8080` | 给程序提供默认监听端口（与代码中读 `PORT` 一致）。 |
| 18 | `ENTRYPOINT ["/app/learn-api"]` | 容器启动时执行的命令；JSON 数组形式**不经 shell**，更直观、少一层解释。 |

**关键：`go build` 那一行（第 8 行）**

- **`CGO_ENABLED=0`**：关闭 cgo，通常得到**静态链接**二进制，在极简 Alpine 里不依赖系统动态库，拷贝到运行时镜像更稳。
- **`-trimpath`**：去掉编译路径前缀，二进制里少带本机路径信息。
- **`-ldflags="-s -w"`**：`-s` 去掉符号表，`-w` 去掉 DWARF 调试信息，**减小可执行文件体积**（调试不便，适合发布镜像）。

**关键：`COPY go.mod` → `go mod download` → 再 `COPY .`**

- 依赖不变时，Docker 可复用「下载依赖」这一层；**只改业务代码**时不必反复 `go mod download`，构建更快。

**关键：`COPY --from=build`**

- **多阶段**的核心：最终镜像里只有 **Alpine + 证书 + 一个二进制 + 非 root**，没有 `go` 命令和完整源码树。

**步骤（在项目根目录或 `learn-api` 下）：**

```bash
cd learn-api
docker build -t learn-api:local -f Dockerfile .
docker run --rm -p 8080:8080 learn-api:local
curl -s localhost:8080/healthz
```

**若 `go mod download` 报 `proxy.golang.org` / `i/o timeout`**：构建容器内默认访问 Google 的 Go 模块代理与校验库，部分网络下会超时。可在 **`learn-api` 目录**用下面命令构建（与 `Dockerfile` 中 `ARG`/`ENV` 对应）：

```bash
docker build -t learn-api:local -f Dockerfile \
  --build-arg GOPROXY=https://goproxy.cn,direct \
  --build-arg GOSUMDB=sum.golang.google.cn \
  .
```

也可换 `https://goproxy.io,direct` 等镜像；能直连官方代理时不必加参数。

**`docker build` 说明**：`-t learn-api:local` 给镜像起名并打标签（`仓库名:标签`，`local` 表示本机构建）；`-f Dockerfile` 指定 Dockerfile；最后的 **`.`** 是**构建上下文**路径，`COPY` 只能引用该目录（及子目录）里的文件。

**kind 加载本地镜像**（给集群用之前必读）：

本机 `docker build` 的镜像在 **宿主机 Docker** 里；kind 集群的 **Node** 是独立容器环境，**默认看不到**你刚打的镜像。若 Deployment 里写 `image: learn-api:local`，不先导入会出现拉取失败。`kind load docker-image` 的作用是把该镜像 **导入当前 kind 集群的各个 Node**（单节点就只导到那一台），效果上接近「推到仓库再在节点上 pull」，适合本地开发省去仓库。

```bash
kind load docker-image learn-api:local --name learn
```

| 部分 | 含义 |
|------|------|
| `kind load docker-image` | 从本机 Docker 把镜像加载进指定 kind 集群的节点。 |
| `learn-api:local` | 已存在的镜像名与标签（须先 `docker build` 成功）。 |
| `--name learn` | 集群名，与 `kind create cluster --name learn` 一致；省略则指向默认集群（常为 `kind`）。 |

若镜像已推到私有仓库且清单里写完整镜像地址，一般**不需要**再 `load`。

**minikube**（与上同理：把本机镜像导入 minikube 节点）：

```bash
minikube image load learn-api:local
```

**练习 2.1**：容器内 `/healthz` 返回 `ok`；本地镜像已加载到集群节点（kind/minikube 按上表执行）。

---

### 2.2 .dockerignore

**要点**：缩小构建上下文、加快 `docker build`，避免把 `.git`、无关文件打进 context。

**练习 2.2**：对比 `docker build` 时 “Sending build context” 大小；在 `learn-api/.dockerignore` 中增加大目录忽略后再次构建，体积应减小或构建更快。
