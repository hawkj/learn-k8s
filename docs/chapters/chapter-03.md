# 第 3 章 Kubernetes 核心对象

**上一章**：[第 2 章](./chapter-02.md)｜**下一章**：[第 4 章 可观测与排障](./chapter-04.md)

---

### 3.1 Deployment

**详细说明（建议第一次读 3.1 时完整过一遍；再动手跑下面命令）**

**1）副本（replicas）是什么**

- **副本**在这里指：**同一套业务逻辑，同时跑几份独立的 Pod**。本仓库 `deployment-api.yaml` 里 **`replicas: 2`**，即希望**始终有 2 个**带标签 `app=api` 的 Pod 在跑（高可用：挂掉一个，另一个仍服务；也可分担流量，视 Service 与副本数而定）。
- **Deployment 不直接等于「一个容器」**：它管的是 **Pod 模板 + 期望副本数**。每个 Pod 里是你在清单里写的 **容器列表**（本示例主要是单个 `api` 容器）。
- **控制器在做什么**：API Server 里保存了「期望状态」；**Deployment 控制器**发现实际 Pod 数量少了（崩溃、节点故障、被删），会**新建** Pod；多了会删；与 `replicas` 不一致就持续调和，所以叫**声明式**管理。

**2）滚动更新（rolling update）是什么**

- **场景**：你改了 Deployment 里的内容（常见是 **`image` 换成新版本**，或改环境变量、资源等），希望**尽量不停机**地切到新版本。
- **做法（默认策略）**：在一段时间内**先起一部分新 Pod**，再**停一部分旧 Pod**，交替进行，直到**全是新 Pod**；而不是「先全删旧再全起新」（那样会断服）。这一过程就叫 **滚动更新**。
- **和 `kubectl apply` 的关系**：你对清单改完再 **`kubectl apply -f deployment-api.yaml`**，Deployment 会检测到模板变化，按策略触发滚动；可用 `kubectl rollout status deployment/api -n demo` 看进度，用 `kubectl rollout undo deployment/api -n demo` 回滚上一版本（需满足版本历史等条件）。
- **进阶**：策略里还有 **`maxSurge` / `maxUnavailable`**（一次最多多几个、最少几个不可用）等，本教程不展开，知道「默认是渐进替换」即可。

**3）「集群」在这里指什么？为什么又提集群？**

- **集群（cluster）**就是你在第 0 章用 **kind / minikube 搭起来的那一套 Kubernetes**：有 API Server、有 Node、有调度器，**和「你本机上的 Docker」不是同一个进程空间**。
- 前面章节你在 **Mac 上** `docker build` 出 **`learn-api:local`**，镜像在 **本机 Docker** 里；而 **`kubectl apply` 是把清单交给「集群里的 API Server」**，由集群在 **Node** 上起 Pod。
- **起 Pod 时**，该 Node 上的 **kubelet** 要能在本节点的**容器运行时镜像库**里找到 `image: learn-api:local`：要么从 **镜像仓库 pull**，要么像 kind 那样 **`kind load` 进节点**。所以文档说 **「镜像需集群可拉取」**——准确说是：**对「将要运行该 Pod 的 Node」而言，镜像必须「已存在或可拉取」**；口语里常说成「集群能拉到镜像」。
- **小结**：「集群」没有新概念，就是你一直在用的 **kind-learn** 那套环境；**「可拉取」**强调 **Mac 上有镜像 ≠ 集群 Node 上一定有**，所以要 **`kind load` / `minikube image load`**（或推仓库再在清单里写仓库地址）。

**要点（摘要）**：`Deployment` = **Pod 模板 + `replicas` 期望副本数**，控制器持续把实际 Pod 数量对齐到该期望；改清单（尤其 `image`）时默认 **滚动更新**。镜像须能在**将运行 Pod 的 Node** 上获取（见上文「3）」）。

**清单原文与字段说明**（与仓库 `deploy/k8s/demo/` 下文件一致；若与 Git 有出入，**以仓库文件为准**。）

#### `namespace.yaml`

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: demo
```

| 字段 | 含义 |
|------|------|
| `apiVersion: v1` | 使用 Kubernetes **核心组** `v1` API（Namespace 属于核心资源）。 |
| `kind: Namespace` | 资源类型为**命名空间**，用于隔离一组资源的名字与配额等。 |
| `metadata.name` | 命名空间名称，此处为 **`demo`**；后续清单里 `namespace: demo`、命令里 `-n demo` 都指向它。 |

#### `configmap-app.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: learn-config
  namespace: demo
data:
  LOG_LEVEL: "info"
```

| 字段 | 含义 |
|------|------|
| `kind: ConfigMap` | 存放**非敏感**键值配置，可被 Pod 以环境变量或文件等形式挂载。 |
| `metadata.name` | 对象名 **`learn-config`**；谁在引用见下。 |
| `metadata.namespace` | 资源所在命名空间 **`demo`**。 |
| `data` | 键值对；键 **`LOG_LEVEL`** 的值 **`info`** 会注入到容器环境变量（由 Deployment 里 `valueFrom` 指定）。 |

**用语**：**`Deployment`** 指本节的 **`deployment-api.yaml` 定义的那种工作负载资源**（`kind: Deployment`），它里面的 **`spec.template.spec.containers[].env`** 描述 Pod 里容器的环境变量。其中某一项可以写 **`valueFrom.configMapKeyRef`**：表示「这个环境变量的值**不要写死在 YAML 里**，而去读某个 **ConfigMap** 里某一键」。**`configMapKeyRef`** 是这段声明的字段名；其下的 **`name`** 对应 ConfigMap 的 **`metadata.name`**（此处为 **`learn-config`**），**`key`** 对应 `data` 里的键（此处为 **`LOG_LEVEL`**）。Secret 的 **`secretKeyRef`** 同理，只是数据源换成 **Secret**。

#### `secret-app.yaml`

```yaml
# 示例：勿提交真实密钥；生产用 sealed-secrets 或外部 Secret 管理
apiVersion: v1
kind: Secret
metadata:
  name: learn-secret
  namespace: demo
type: Opaque
stringData:
  API_KEY: "demo-local-key"
```

| 字段 | 含义 |
|------|------|
| 文件首行注释 | 提醒：**不要**把真实密钥写进 Git；生产可用 SealedSecrets、Vault 等方案。 |
| `kind: Secret` | 存放**敏感**数据（密码、Token 等）；落库时会被编码，访问受 RBAC 约束。 |
| `metadata.name` / `namespace` | 对象名 **`learn-secret`**，位于 **`demo`**。 |
| `type: Opaque` | 通用自定义类型；用户任意键值，由 Kubernetes 按 Secret 规则存储。 |
| `stringData` | 以**明文**写在 YAML 里，**apply 时**由 API 转为内部存储形式（便于本地示例）；另一种写法是 `data` 里放 **base64** 字符串。 |
| `stringData.API_KEY` | 键名 **`API_KEY`**，供 Deployment 里 `secretKeyRef.key` 引用；值 **`demo-local-key`** 仅作本地演示。 |

#### `deployment-api.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: demo
  labels:
    app: api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - name: api
          # kind: docker build -t learn-api:local learn-api && kind load docker-image learn-api:local
          image: learn-api:local
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: PORT
              value: "8080"
            - name: LOG_LEVEL
              valueFrom:
                configMapKeyRef:
                  name: learn-config
                  key: LOG_LEVEL
            - name: API_KEY
              valueFrom:
                secretKeyRef:
                  name: learn-secret
                  key: API_KEY
            # 第 6 章启用链式调用时取消注释：
            # - name: BACKEND_BASE_URL
            #   value: "http://backend.demo.svc.cluster.local:8080"
          resources:
            requests:
              cpu: "50m"
              memory: "64Mi"
            limits:
              cpu: "500m"
              memory: "256Mi"
          livenessProbe:
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 3
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /readyz
              port: http
            initialDelaySeconds: 1
            periodSeconds: 5
```

**Deployment / metadata**

| 字段 | 含义 |
|------|------|
| `apiVersion: apps/v1` | **工作负载** API 版本；Deployment 属于 `apps` 组。 |
| `kind: Deployment` | 控制器资源：根据模板维护多份 Pod，并支持滚动更新等。 |
| `metadata.name` | Deployment 对象名 **`api`**（`kubectl get deployment -n demo api`）。 |
| `metadata.namespace` | 部署在 **`demo`**。 |
| `metadata.labels` | Deployment **自身**的标签（可用于筛选、给别的资源引用等）；与下面 Pod 标签常一致但不强制相同概念。 |

**`spec`（期望状态）**

| 字段 | 含义 |
|------|------|
| `replicas: 2` | 期望 **2** 个符合模板的 Pod 同时运行。 |
| `selector.matchLabels` | Deployment **认领**哪些 Pod：必须带 **`app: api`** 标签，且与 `template.metadata.labels` **一致**，否则 apply 会失败。 |
| `template` | **Pod 模板**：每新建一个 Pod，就按这里描述来生成。 |
| `template.metadata.labels` | 每个 Pod 上的标签；须满足上面的 **`selector`**。 |

**`template.spec.containers`（本例仅一个容器）**

| 字段 | 含义 |
|------|------|
| `name: api` | 容器名（多容器时区分；`kubectl logs -c api` 会用到）。 |
| `image` | 使用的镜像；本地 kind 需 **`kind load`** 后节点才有 `learn-api:local`。 |
| `imagePullPolicy: IfNotPresent` | 节点上**已有**该镜像则不向仓库拉取；没有才拉取（本地开发常用）。 |
| `ports.containerPort` | 容器内进程监听端口 **8080**。 |
| `ports.name: http` | 端口名，探针里 **`port: http`** 通过名字引用该端口。 |
| `env`（`PORT`） | 直接写死的普通环境变量。 |
| `env`（`LOG_LEVEL`） | **`valueFrom.configMapKeyRef`**：从 ConfigMap **`learn-config`** 的键 **`LOG_LEVEL`** 注入。 |
| `env`（`API_KEY`） | **`valueFrom.secretKeyRef`**：从 Secret **`learn-secret`** 的键 **`API_KEY`** 注入。 |
| 注释掉的 `BACKEND_BASE_URL` | 第 6 章再启用；指向集群内 Service DNS。 |
| `resources.requests` | **调度**与 **QoS** 参考：容器至少声明需要这么多 CPU/内存（`m` 表示 millicores，`Mi` 表示 Mebibytes）。 |
| `resources.limits` | **上限**：超过可能被 **OOMKill** 或 CPU **节流**（见 3.6）。 |
| `livenessProbe` | **存活探针**：失败则**重启容器**；HTTP 访问 **`/healthz`**，端口名 **`http`**。 |
| `livenessProbe.initialDelaySeconds` | 容器启动后等待 **3** 秒再开始探测。 |
| `livenessProbe.periodSeconds` | 每 **10** 秒探测一次。 |
| `readinessProbe` | **就绪探针**：失败则**从 Service 摘流**（不重启）；路径 **`/readyz`**。 |
| `readinessProbe.initialDelaySeconds` / `periodSeconds` | 就绪探测的首次延迟与周期。 |

**本节步骤的整体目的**（和「设置集群」的区别）

- **不是在安装或配置整个 Kubernetes 集群**。集群本身应已用 **kind / minikube** 建好并能 `kubectl get nodes`；那一步才是「有集群可用」。
- **整体目的**：在**已有集群里**，用清单**声明**一套 **demo 里的 API 演示环境**——先有命名空间和配置/密钥，再让 **Deployment** 按模板起 **2 个 `api` Pod**；最后一条 **`kubectl get pods`** 用来**验收** Pod 是否已创建并进入预期状态。
- 可记成：**`apply` = 把「期望长什么样」交给控制面；控制面在集群里替你落实**（起 Pod、挂配置等）。

**步骤：**

```bash
kubectl apply -f deploy/k8s/demo/namespace.yaml
kubectl apply -f deploy/k8s/demo/configmap-app.yaml
kubectl apply -f deploy/k8s/demo/secret-app.yaml
kubectl apply -f deploy/k8s/demo/deployment-api.yaml
kubectl get pods -n demo -l app=api
```

**命令说明**（均在仓库根目录执行；`kubectl apply` **不是**在「每次执行时固定只做创建或只做更新」里二选一，而是：**看集群里是否已有该资源**——**没有则创建，有则按清单做声明式对齐**（有字段变化会更新/打补丁，与清单完全一致且策略允许时也可能显示 **`unchanged`**）。）

| 命令 | 作用 |
|------|------|
| `kubectl apply -f …/namespace.yaml` | 若集群里**还没有** `demo` 命名空间 → **创建**；**已有**且与 `namespace.yaml` 一致 → 通常 **`unchanged`**；已有但你改过清单 → **更新**（Namespace 可改字段很少，多数学习场景是第一次创建、之后反复 apply 多为 unchanged）。后续资源都加 **`-n demo`**。 |
| `kubectl apply -f …/configmap-app.yaml` | 创建或更新 **ConfigMap**（非敏感配置，如日志级别）；Deployment 里通过 `envFrom` / `volume` 等引用（以清单为准）。 |
| `kubectl apply -f …/secret-app.yaml` | 创建或更新 **Secret**（敏感或需 base64 的字段）；示例仅作本地学习，生产勿把真密钥提交 Git。 |
| `kubectl apply -f …/deployment-api.yaml` | 创建或更新 **Deployment `api`**：声明副本数、镜像、环境变量来源等；控制器会据此创建/更新 **Pod**。 |
| `kubectl get pods -n demo -l app=api` | 列出 **`demo` 命名空间**里、带标签 **`app=api`** 的 Pod（与 Deployment 里 Pod 模板标签对应）；用于确认是否已调度、`Running` 等。 |

**练习 3.1**：`demo` 命名空间下 **2** 个 `api` Pod 均为 `Running`。

---

### 3.2 Service（ClusterIP）

**要点**：`Service` 通过 **selector** 匹配 Pod，提供稳定 **DNS**：`api.demo.svc.cluster.local`（同命名空间内可简写 `api`）。

**说明：selector 是什么？为什么要「稳定 DNS」？**

- **`selector`（选择算子）**：写在 **`Service` 的 `spec.selector`** 里，是一组 **标签键值对**（本仓库 `service-api.yaml` 里是 **`app: api`**）。集群会把 **带相同标签的 Pod** 当作这个 Service 的**后端**；Deployment 里 Pod 模板的 **`metadata.labels`** 必须与之**一致**，流量才会打到你的 `api` Pod 上。Pod 重建、扩缩容时，只要标签还在，**Service 会自动把新 Pod 纳入或摘掉**，你不必改 Service YAML。
- **为什么要稳定 DNS**：**Pod 的名字、IP 会随重建而变化**（例如滚动更新后变成新 Pod 名、新 IP）。若业务写死「连某个 Pod IP」，一改就断。**`Service` 自己有固定名字**（此处 **`api`**），集群 **DNS**（如 **CoreDNS**）把名字解析到「当前符合 selector 的一组 Pod 地址」。因此客户端只要访问 **`http://api:8080`**（同命名空间短名）或全名 **`api.demo.svc.cluster.local`**，就能在 Pod 换来换去时仍找到**可用的后端**——这就是「稳定」的含义：**稳定的是 Service 名与 DNS，不是某个 Pod 实例**。

**`service-api.yaml` 原文与字段**（`deploy/k8s/demo/service-api.yaml`；与 Git 不一致时以仓库为准。）

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api
  namespace: demo
spec:
  selector:
    app: api
  ports:
    - name: http
      port: 8080
      targetPort: http
```

| 字段 | 含义 |
|------|------|
| `apiVersion: v1` | Service 属于核心 API **`v1`**。 |
| `kind: Service` | 资源类型为 **Service**（集群内虚拟 IP + DNS + 负载均衡到后端 Pod）。 |
| `metadata.name` | Service 名 **`api`**；同命名空间内 DNS 短名即为 **`api`**。 |
| `metadata.namespace` | 位于 **`demo`**；全名 DNS 为 **`api.demo.svc.cluster.local`**。 |
| `spec.type`（本文件未写） | 省略时默认为 **`ClusterIP`**：仅在集群内可访问，不自动暴露到宿主机（与 **NodePort / LoadBalancer** 相对）。 |
| `spec.selector` | 只把带 **`app: api`** 标签的 Pod 作为后端；须与 **Deployment** 里 Pod 模板标签一致。 |
| `spec.ports[].name` | 端口记录名 **`http`**，可与 **`targetPort`** 同名，便于探针、其它资源引用端口名。 |
| `spec.ports[].port` | **Service 对外提供的端口**（ClusterIP 上监听 **8080**）；集群内访问 **`api:8080`** 即连到此端口。 |
| `spec.ports[].targetPort` | 流量**转发到 Pod 里容器的哪个端口**；写 **`http`** 表示用**名字**对齐 Pod 里 **`containerPort` 的 `name: http`**（见 `deployment-api.yaml`）；也可写数字如 **`8080`**。 |

**步骤：**

```bash
kubectl apply -f deploy/k8s/demo/service-api.yaml
# 若集群拉 Hub 不稳定，可先在本机 docker pull 再 kind load，并给镜像加固定 tag（避免 :latest 默认 Always 拉取）：
# docker pull curlimages/curl:8.5.0 && kind load docker-image curlimages/curl:8.5.0 --name learn
kubectl run curl --rm -it --restart=Never -n demo --image=curlimages/curl:8.5.0 --pod-running-timeout=5m -- \
  curl -s http://api:8080/healthz
```

**练习 3.2**：集群内访问 `http://api:8080/healthz` 返回 **`ok`**（通过上条 `kubectl run … curl` 验证）。

**可跳过练习吗**：可以。**练习 3.2** 的目的只是确认 **Service `api` 在集群内可达**；若你已 **`apply` `service-api.yaml`**，且 **`kubectl get pods -n demo -l app=api`** 均为 **`Running`**、**`kubectl get endpoints -n demo api`** 里能看到 **Pod IP**，则 **Service + DNS + 后端** 在概念上已对齐，**不必强求本机一定跑通临时 curl Pod**；可先继续 **3.3** 及后文，网络或代理（见下）就绪后再回来补做。

**若 `kubectl run …` 一直卡住、或 Pod 为 `ImagePullBackOff`**

- **`kubectl get pods` 只显示状态**（如 `ImagePullBackOff`），**不会**打出完整拉镜像错误；详细原因在 **`kubectl describe pod <Pod名> -n demo`** 末尾的 **`Events:`** 里（`Failed` / `Pulling` 的 **Message**）。
- 也可：`kubectl get events -n demo --field-selector involvedObject.name=curl --sort-by='.lastTimestamp'`
- **常见（kind / Colima）**：事件里出现 **`proxyconnect tcp: dial tcp 127.0.0.1:7897: connect: connection refused`**，表示节点拉镜像时被配置成走 **`http://127.0.0.1:7897` 代理**，但该地址上**没有服务**（本机代理未开，或节点里的 `127.0.0.1` 并非你 Mac 上的代理）。处理思路：打开本机代理并修正 kind/虚拟机内代理地址、或建集群时勿向节点注入失效的 `HTTP_PROXY` / 重建集群等；再配合上文 **`docker pull` + `kind load` + 固定 tag`** 减少对外拉取次数。

---

### 3.3 ConfigMap / Secret

**本节目的**：学会把**可变的配置**从容器镜像里拆出来，交给 Kubernetes 的 **ConfigMap**（一般配置）和 **Secret**（敏感数据）管理；Deployment 只**引用**它们，这样改日志级别、换密钥等**不必重新 build 镜像**，改 YAML `apply` 并（必要时）重启 Pod 即可。同时建立习惯：**真密钥不进 Git**，示例 Secret 仅供本地。

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
