# Startup Protocol

## Purpose

定义进入当前仓库后的默认 agent 启动方式。

## Scope

适用于所有非简单需求。

## Must

- 读取顺序：
  - `AGENTS.md`
  - `agents/README.md`
  - `agents/manifest.yaml`
  - `agents/protocols/startup.md`
  - `agents/roles/*.md`
  - `agents/reference/worker-constraints.md`
  - `agents/reference/evaluator-rubric.md`
  - `agents/reference/observability.md`
- 默认加载 3 个角色：
  - `Planner`
  - `Generator`
  - `Evaluator`
- 运行时支持 delegation 时，必须真实创建 3 个 sub-agent，而不是在单 agent 内模拟角色
- 非简单需求先初始化 `agents/runs/YYYYMMDD-<slug>/`

## Must Not

- 不得跳过 run 目录
- 不得没有 contract 就直接实现
- 不得没有 evaluator report 就直接宣布完成
- 不得省略 sub-agent 协作约束

## Outputs

- 初始化后的 run 目录
- 后续各角色工件

## Pass/Fail

- `Pass`：3 角色已进入同一个 run 上下文
- `Fail`：角色、工件、协议三者任一缺失
