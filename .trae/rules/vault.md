# 项目规范（go-mpc-vault）

## 总则
- 统一使用 `Go 1.24.x`（模块声明见 `go.mod`）。
- 所有开发任务通过 `Makefile` 驱动；避免直接调用底层工具命令。
- `internal/models` 与 `internal/types` 目录只存放生成代码；不得手写。
- 脚本代码需设置构建标签 `//go:build scripts`。
- 提交信息遵循 Conventional Commits；PR 遵循统一检查清单。

## 目录结构
- 根目录关键文件
  - `go.mod`（`github.com/kashguard/go-mpc-vault`）
  - `Makefile`（构建、测试、SQL、生成工具）
  - `Dockerfile`、`docker-compose.yml`（开发容器与依赖服务）
  - `.golangci.yml`（Lint 配置）、`.drone.yml`（CI）
  - `internal/`（业务与基础设施代码）
  - `migrations/`（数据库迁移）
  - `api/`（HTTP 层、OpenAPI/Swagger 资源）
  - `scripts/`（带 `//go:build scripts` 的工具或脚本）

## 开发流程
- 初始化
  - `make init`：下载依赖、安装工具、`go mod tidy`
  - `make modules`：缓存 go modules
  - `make tools`：安装工具（见 `go.mod` tool 块）
- 生成与预处理
  - `make build-pre`：执行 SQL、Swagger、生成绑定、`go generate`
  - `make sql`：格式化 SQL、重置并迁移 Spec DB、`sqlboiler` 生成模型
  - `make go-generate-handlers`：生成 `internal/api/handlers/handlers.go` 绑定
- 构建与格式化
  - `make go-format`：`go fmt ./...`
  - `make go-build`：`go build -ldflags $(LDFLAGS) -o bin/app`
- Lint
  - `make lint`：`golangci-lint run --timeout 5m`
- 测试
  - `make test`：包维度输出、生成覆盖率文件 `/tmp/coverage.out`
  - `make test-by-name`：用例维度输出
  - `make go-test-print-coverage`：打印总体覆盖率
  - `make go-test-print-slowest`：打印最慢用例（阈值 2s）
  - `make test-update-golden`：刷新黄金文件（需二次确认）
- 信息与检查
  - `make info`：打印 DB/Handlers/go.mod 变更信息与当前 Go 版本
  - `make check-gen-dirs`：确保 `internal/models|types` 仅含生成文件
  - `make check-script-dir`：确保 `scripts/*.go` 含 `//go:build scripts`

## SQL 与数据模型
- 迁移文件位于 `migrations/`；修改后需运行 `make sql` 以生成模型。
- SQL 格式化使用 `pg_format`；语法检查通过真实数据库（见 `make sql-check-syntax`）。
- `sqlboiler psql` 基于 Spec DB 生成 `internal/models/*.go`。
- 禁止手工编辑生成的模型文件；如需扩展，使用 `repository` 或服务层封装。

## HTTP 与 API
- HTTP 框架：Echo v4（`github.com/labstack/echo/v4`）。
- 处理器绑定文件由生成工具维护（`gsdev handlers gen`）；不要手改。
- 参数校验、DTO 定义放置在 `internal/api` 与 `internal/model`。
- OpenAPI/Swagger 文档需与实现保持一致；规范变更需更新生成。

## 配置与日志
- 配置管理：Viper；环境变量示例见 `.env.local.sample`。
- 日志：Zerolog（结构化日志）；使用 request-scoped logger 与 `context.Context`。
- 错误：优先 `github.com/pkg/errors` 包装并保留调用栈；统一返回码与错误映射。

## 容器与依赖
- 使用 `docker-compose.yml` 启动 Postgres、Integresql 与其他依赖。
- 数据库默认开发配置：`PGHOST=postgres`、`PGDATABASE=mpc-dev-db`。
- 开发容器内运行构建与测试，统一环境与架构。

## Make 命令执行约定
- 所有 `make` 命令必须在开发容器内执行；在宿主机运行会因 `SHELL=/app/rksh` 等环境依赖失败。
- 进入方式
  - VS Code Dev Containers：在容器终端（如 `development@...:/app`）直接运行 `make ...`
  - Docker Compose：先运行 `docker compose up -d`，然后执行 `docker compose exec service make <target>`
- 示例
  - `docker compose exec service make init`
  - `docker compose exec service make build`
  - `docker compose exec service make test`
- 环境变量示例
  - `docker compose exec -e GO_MODULE_NAME=github.com/kashguard/go-mpc-vault service make get-module-name`

## 安全与合规
- 不得在仓库提交敏感信息（密钥、密码、CA 私钥等）。
- 对外接口的关键操作需要二次验证（Passkey/WebAuthn）；规范参考上游文档。
- 所有 HTTP/后台任务需带超时与重试；幂等性设计要求明确（特别是 DKG/签名触发）。

## 分支与提交规范
- 分支策略：Trunk-Based Development；短分支、快合并。
- Commit 消息：Conventional Commits
  - `feat: ...` 新功能
  - `fix: ...` 修复
  - `docs: ...` 文档
  - `refactor: ...` 重构
  - `test: ...` 测试
  - `chore: ...` 杂项（依赖升级、脚手架等）
- PR 检查清单
  - 通过 `make lint` 与 `make test`
  - 更新/生成必要的绑定与模型
  - 无生成目录的手改代码（`internal/models|types`）
  - 文档与配置同步更新

## 版本与发布
- 变更日志参考 `CHANGELOG.md` 与 `CHANGELOG-go-starter.md`。
- CI/CD 管道见 `.drone.yml`；发布分支与制品命名遵循 CI 配置要求。

## 常用命令速查
- 以下命令均需在容器内执行（见上节）。
- 初始化：`make init`
- 生成模型：`make sql`
- 构建：`make build` 或 `make go-build`
- Lint：`make lint`
- 测试：`make test`
- 打印覆盖率：`make go-test-print-coverage`
- 检查生成目录：`make check-gen-dirs`
- 检查脚本标签：`make check-script-dir`
