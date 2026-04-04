# Harness 规范

## 默认流程

本仓库默认使用 `Planner -> Generator -> Evaluator` 三段式。

## 执行要求

- 运行时支持 delegation 时，必须真实启动 3 个 sub-agent
- 不支持 delegation 时，也要保留 3 个角色的 contract、handoff 和 evaluator 结论
- 非简单需求必须创建 `agents/runs/YYYYMMDD-<slug>/`

## 角色职责

- `Planner`：拆分任务、定义边界、写清验收标准
- `Generator`：实施变更
- `Evaluator`：检查边界、风险和测试覆盖

## 输出工件

- `request.md`
- `product-spec.md`
- `sprint-01-contract.md`
- `handoff.md`
- `evaluator-report.md`

## 适用范围

- verify 接口变更
- 写回逻辑调整
- 观测规则更新
- Swagger 同步
