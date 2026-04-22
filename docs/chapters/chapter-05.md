# 第 5 章 CI/CD

**上一章**：[第 4 章](./chapter-04.md)｜**下一章**：[第 6 章 多服务](./chapter-06.md)

---

### 5.1 流水线阶段

**要点**：典型顺序：**lint** → **单元测试** → **构建二进制** → **构建镜像** → **推送仓库** → **更新集群**。

**本仓库**：`.github/workflows/ci.yml` 在 `learn-api` 下执行 `go test ./...`。

**`.github/workflows/ci.yml` 逐段说明**

文件路径固定为 **`.github/workflows/*.yml`**，GitHub 才会把它当作 **GitHub Actions** 工作流。下面按出现顺序解释（与仓库当前文件一致）。

```yaml
name: ci

on:
  push:
    branches: [main, master]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: learn-api
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - run: go test ./...
```

| 片段 | 含义 |
|------|------|
| **`name: ci`** | 在 GitHub 仓库 **Actions** 页里显示的工作流名称，便于识别。 |
| **`on:`** | **何时触发**本流水线。 |
| **`push: branches: [main, master]`** | 有人向 **`main` 或 `master` 分支推送提交**时运行（含直接 push、merge 进主分支等）。 |
| **`pull_request:`** | 有人**打开或更新 Pull Request** 时运行（用于在合并前跑检查）；未限定分支即对指向本仓库的常见 PR 事件生效。 |
| **`jobs:`** | 定义一个或多个 **Job**（任务组）；彼此可设依赖，本文件只有一个。 |
| **`test:`** | Job 的 **ID**（内部名），在 Actions 界面里常显示为 **`test`**。 |
| **`runs-on: ubuntu-latest`** | 在 GitHub 托管的 **Ubuntu 最新版虚拟机**上执行；每次运行通常是干净环境。 |
| **`defaults.run.working-directory: learn-api`** | 后面所有 **`run:` shell 步骤**的**当前目录**设为仓库里的 **`learn-api/`**，这样 **`go test ./...`** 在正确模块路径下执行（与本地 `cd learn-api && go test ./...` 一致）。 |
| **`steps:`** | Job 内的**有序步骤**，上一步失败则后续一般不再跑。 |
| **`- uses: actions/checkout@v4`** | 使用官方 **checkout** 动作：把**当前触发事件对应的代码**检出到 Runner 工作区（浅克隆等由该 action 处理）。 |
| **`- uses: actions/setup-go@v5`** | 安装并缓存 **Go 工具链**；`with.go-version` 指定版本，需与 `go.mod` / 教程一致（此处 **1.22**）。 |
| **`- run: go test ./...`** | 在 **`learn-api`** 目录执行 **Go 测试**：**`./...`** 表示**当前模块下所有包**（递归子目录），与第 1 章本地跑测试一致。 |

**和「CD」的关系**：本工作流只做 **持续集成（CI）** 里的**测试一环**；**构建镜像、推仓库、`kubectl` 部署**等属于更完整的流水线（见 5.2、5.3），可在此基础上往 `jobs` 里加 step 或加新 job。

**练习 5.1**：将仓库推送到 GitHub，确认 Actions 中 **test** 任务通过（无 GitHub 时本地执行 `cd learn-api && go test ./...` 等价）。

---

### 5.2 镜像推送

**要点**：本教程 **CI 只使用 GHCR**（`ghcr.io`）：在 job 里用 **`docker/login-action`** 登录，再 **`docker build` + `docker push`**；镜像 tag 建议用 **`git` 提交短 SHA**。**Docker Hub 为可选方案**，需要时自行改 `registry` 与 Secrets，正文不再展开。

**本节目的**

- **在流水线里产出「可部署」的镜像**：与第 2 章本机 `docker build` 类似，但由 **GitHub Runner** 执行，保证**谁 push、谁合并**，产物都来自**同一份已检出的代码**，减少「我机器上能跑、别人/集群不行」的差异。  
- **推到 GHCR**：集群上的 kubelet 从 **Registry 拉镜像**；镜像在 **`ghcr.io/...`** 后，第 5.3 的 **`kubectl set image` 或清单改 tag** 才能指向**集群能访问**的地址（本地 kind 的 `kind load` 是另一条路，CI 场景以 GHCR 为主）。

**GHCR 镜像名**

| 项 | 说明 |
|------|------|
| **典型名** | `ghcr.io/<owner小写>/learn-api:<tag>` |
| **登录** | 同一 GitHub 生态下，CI 里用 **`GITHUB_TOKEN`** + **`packages: write`** 即可（见下「如何在 CI 中登录」）。 |

**推荐 tag**：**`<git-sha>`**（如 `github.sha`），一眼对应某次提交；也可用 **`main-<date>`** 等，团队约定一致即可。

**在流水线里增加 build / push：每一步在干什么**

下面按**常见顺序**说明「为什么要做这一步」（与是否使用 `docker/build-push-action` 无关，逻辑相同）。

| 步骤 | 目的 |
|------|------|
| **检出代码**（如 `actions/checkout`） | Runner 默认空目录；要有 **`Dockerfile` 与源码`** 才能 `docker build`；构建上下文通常含 **`learn-api/`**（与第 2 章一致）。 |
| **登录 GHCR**（`docker/login-action`） | **`docker push` 需要认证**；不登录则推不上去。传 **`registry: ghcr.io`**、用户名、密码（见下「如何在 CI 中登录」）。 |
| **`docker build`** | 在 Runner 上生成镜像层，**`-t`** 打成 **`ghcr.io/.../learn-api:<sha>`** 这种**带仓库地址与 tag** 的名字，便于下一步直接 **`push` 同名镜像**。 |
| **`docker push`** | 把本地（Runner 上）镜像**上传到 Registry**；成功后集群或其它环境才能 **`pull`** 同一地址。 |

**鉴权（GHCR + CI）**

- **`GITHUB_TOKEN`**：每个 Workflow 运行时会自动注入，**不必**手抄密钥。要对 **GHCR** 执行 **`docker push`**，必须在**同一 job** 声明 **`permissions: packages: write`**（常与 `contents: read` 一起写）。**目的**：让 CI 有权向 **GHCR** 写入包。  
- **PAT（Personal Access Token，个人访问令牌）**：**本机** `docker login ghcr.io` 时用；或 CI 若不用 `GITHUB_TOKEN`、改用仓库 **Secret** 里存的 PAT 时，在 `docker/login-action` 的 **`password:`** 里引用 **`${{ secrets.XXX }}`**。**不要把 PAT 明文写进 YAML**。

**在 GitHub 上哪里能看到新镜像 / tag**

- **GHCR**：仓库页面 **右侧「Packages」**（或组织 **Packages**），找到 **`learn-api`**（或你推送时用的包名）；点进去可看 **tag / digest**、公开范围等。

**如何在 CI 中用 `docker/login-action` 登录 GHCR（实现步骤）**

1. **放在会执行 `docker push` 的同一个 job 里**（例如下文 **`image`** job），不要单独开一个只 login 却不 push 的 job（无意义）。  
2. **在该 job 顶层**增加 **`permissions`**（与 `runs-on` 同级），至少包含：

   ```yaml
   permissions:
     contents: read
     packages: write
   ```

   **`packages: write`**：允许本 job 使用的 **`GITHUB_TOKEN`** 向 **GHCR** 写入；没有它，`docker push` 常会 **403**。  
3. **在 `steps` 里、`docker build` / `docker push` 之前**，增加一步（版本号 `v3` 可按 [docker/login-action](https://github.com/docker/login-action) 仓库说明更新）：

   ```yaml
   - uses: docker/login-action@v3
     with:
       registry: ghcr.io
       username: ${{ github.actor }}
       password: ${{ secrets.GITHUB_TOKEN }}
   ```

   | 字段 | 含义 |
   |------|------|
   | **`registry: ghcr.io`** | 固定写 **GHCR** 的域名。 |
   | **`username`** | 一般用 **`${{ github.actor }}`**（触发本次运行的 GitHub 用户，多为小写用户名）；与 GHCR 登录约定一致即可。 |
   | **`password`** | 填 **`${{ secrets.GITHUB_TOKEN }}`**：GitHub 为本次运行**临时签发**的 token，**不要**换成你自己键盘输入的字符串。 |

4. **其后**再接 **`run: docker build ...` / `docker push ...`**（或 `docker/build-push-action`）。**顺序必须是**：先 **checkout**（若需要源码）→ **login** → **build** → **push**。

完成以上四步，即实现「**在 CI 中使用 `docker/login-action` 登录 GHCR**」。完整 job 示例见下文 **实施步骤 B**。

**推荐做法：本机 push 与 CI push 如何分工**

| 方式 | 典型用途 | 是否作为团队「正规路径」 |
|------|----------|---------------------------|
| **只在 YAML 里配 CI**（Runner 上 `login` + `build` + `push`） | 合并进 **`main`** / 打 tag 后**自动**产出与提交绑定的镜像 | **是**，多数团队日常依赖这条 |
| **本机 `docker login` + push** | 第一次接 GHCR、验证 PAT、网络/权限**排障**、CI 未建好时的临时推送 | **辅助**，不替代 CI 的长期规范 |

二者可并存：**正式环境用的 tag 以 CI 推的为准**；本机用于试验与对照。

---

**实施步骤 A：本机登录 GHCR 并推送（第一次 / 排障）**

1. **创建 PAT**（仅本机 push 需要）：GitHub → **Settings → Developer settings → Personal access tokens**，新建 **classic** 或 **fine-grained**（按页面说明），至少勾选 **`write:packages`**（拉私有包时常加 **`read:packages`**）。**不要**把 PAT 写进仓库，只保存在本机密码管理或终端临时粘贴。  
2. **登录**：在终端执行（将 **`<GitHub用户名>`** 换成你的用户名；密码处粘贴 **PAT**，输入时不可见为正常）：

   ```bash
   docker login ghcr.io -u <GitHub用户名>
   ```

3. **构建并打 tag**（在仓库根目录；与第 2 章一致，**`-t` 必须含 `ghcr.io/...`**）：

   ```bash
   docker build -t ghcr.io/<GitHub用户名小写>/learn-api:manual-test -f learn-api/Dockerfile learn-api
   ```

   **说明**：GHCR 路径中的 **owner 一般要求全小写**；若用户名有大写，请改成小写或使用你在 **Packages** 里实际看到的包路径。  
4. **推送**：

   ```bash
   docker push ghcr.io/<GitHub用户名小写>/learn-api:manual-test
   ```

5. 到 GitHub **Packages** 确认出现 **tag `manual-test`**（或你起的名字）。

---

**实施步骤 B：在 GitHub Actions 里推送（推荐「正规路径」）**

**目的**：每次 **`test` 通过**且代码 **`push` 到 `main`/`master`** 时，自动构建并推到 GHCR，tag 用 **`github.sha`** 对齐提交。

**落盘方式**：编辑 **`.github/workflows/ci.yml`**（或新建 **`release.yml`** 等），在 **`jobs:`** 下与 **`test`** job **同级**增加 **`image:`** job，并写 **`needs: test`**，使镜像构建**排在单测之后**。

下面为**可整段追加的示例**（追加到现有 `ci.yml` 的 `jobs:` 内；**缩进与 `test` 对齐**）。若你的 Docker 构建上下文与下文不一致，只改 **`docker build`** 那一行的 **`-f`** 与**最后一项上下文路径**即可。

```yaml
  image:
    name: Build and push learn-api
    needs: test
    if: github.event_name == 'push' && (github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master')
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Docker build and push
        run: |
          OWNER_LC=$(echo "${{ github.repository_owner }}" | tr '[:upper:]' '[:lower:]')
          IMAGE="ghcr.io/${OWNER_LC}/learn-api:${{ github.sha }}"
          docker build -t "${IMAGE}" -f learn-api/Dockerfile learn-api
          docker push "${IMAGE}"
```

| 片段 | 目的 |
|------|------|
| **`needs: test`** | **单测通过**后再构建镜像，避免浪费 Registry 与 Runner。 |
| **`if: ... push ... main/master`** | 仅在**推送到主分支**时推镜像；**PR** 只跑 `test`、不推（减少 Fork PR 与密钥风险；与「注意」一节一致）。可按需删掉 `if` 或改成「仅 tag」。 |
| **`permissions.packages: write`** | 允许 **`GITHUB_TOKEN`** 向 **GHCR** 写入（见上文「鉴权」）。 |
| **`OWNER_LC=... tr`** | 将 **`repository_owner`** 转**小写**，减少 GHCR 拒收概率。 |
| **`docker/login-action` + `GITHUB_TOKEN`** | 在 Runner 上完成 **GHCR 登录**，无需把 PAT 写进仓库。 |

首次成功后：仓库 **Packages** 中会出现 **`learn-api`**（或首次 push 时自动创建的包），tag 为 **提交 SHA**。第 **5.3** 将镜像地址改为 **`ghcr.io/<owner>/learn-api:<该 SHA>`** 即可让集群拉取（集群需能访问公网 GHCR，或配 **imagePullSecret**；本教程不展开）。

---

**练习 5.2**：按 **实施步骤 B** 将 **`image` job** 合入工作流（或先完成 **实施步骤 A** 理解登录与 push），在 **Actions** 中看到 **`image`/`Build and push learn-api`** 成功，并在 **Packages** 中看到 **新 tag**。若仅合并文档未改 `ci.yml`，可只完成 **A** 作为手工演练。

**注意**

- **Fork 提 PR**：默认 **`GITHUB_TOKEN`** 对**上游仓库包**可能没有写权限，**push 镜像到上游 GHCR** 常会失败；通常只在 **`push` 到本仓库主分支** 时推镜像，或改用有权限的 token / 单独发布流程。  
- **GHCR 包名**：部分场景要求 **全小写**，若 push 被拒可检查 **owner / 镜像名** 是否与 **Packages** 要求一致。

---

### 5.3 部署到集群

**要点**：

- **方式 A**：`kubectl set image deployment/api api=<新镜像> -n demo`
- **方式 B**：Kustomize `images.newTag` 或 Helm `values`
- CI 需 **kubeconfig** 以 **Secret** 注入（注意安全）

**练习 5.3**：在 CI 或本机用可访问集群的 kubeconfig 执行一次 **apply** 或 **set image**，使 `api` 使用新镜像并成功滚动（本地 kind 可先手动 `kubectl set image` 验证）。

---

### 5.4 回滚

**要点**：`kubectl rollout history deployment/api -n demo`、`kubectl rollout undo deployment/api -n demo`。

**练习 5.4**：将 Deployment 镜像改为 `busybox` 等无法启动的应用，待 Pod 异常后执行 `rollout undo`，恢复上一版。
