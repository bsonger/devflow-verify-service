# Devflow Verify Service

`devflow-verify-service` 只负责 `Verify` 入站接口，以及 verify 写回所需的最小 `Manifest`、`Release`、`Intent` 内部能力。

边界：

- 仅暴露 `/api/v1/verify/*`
- 仅保留 verify 写回需要的最小 model、service、store、配置加载
- 不提供 `Project`、`Application`、`Configuration`、`Manifest`、`Release`、`Intent` 的对外 CRUD 面
- 启动、路由中间件和观测基础设施来自 `../devflow-service-common`

仓库文档：

- [架构](docs/architecture.md)
- [接口规范](docs/api-spec.md)
- [约束](docs/constraints.md)
- [观测规范](docs/observability.md)
- [Harness](docs/harness.md)
- [资源说明](docs/resources/README.md)

运行约定：

- 任何调用其他服务或外部系统的代码都必须同时产出 `metrics + trace + structured log`
- 默认 harness 为 `Planner -> Generator -> Evaluator`
- 运行时支持 delegation 时，必须真实启动对应 sub-agent，不允许只在单 agent 内口头模拟

常用命令：

- `go run ./cmd`
- `go build ./cmd/main.go`
- `go test ./...`
- `swag init -g cmd/main.go --parseDependency`
