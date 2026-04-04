# 约束

## 服务边界

- 本仓库对外只允许 `Verify`
- `Manifest`、`Job`、`Intent` 仅作为 verify 写回的内部依赖
- 不允许重新引入其他资源的对外 router、handler、Swagger 面

## 写回约束

- 写回必须按资源 ID 精确更新
- 对已终态 step 的重复回写应保持幂等
- `pipeline_id`、`task_name`、`step_name` 不能为空字符串
- token 校验必须在写接口生效

## 观测约束

- 任何调用其他服务或外部系统的代码都必须同时产出 `metrics + trace + structured log`
- 不允许把 `job_id`、`manifest_id`、`intent_id`、`external_ref` 作为 metrics label
- 这些业务主键只能进入日志字段和 trace attributes

## 文档约束

- `README.md`、`AGENTS.md`、`agents/protocols/startup.md`、`docs/*.md` 必须描述 verify-only 边界
- Swagger 只允许出现 `/api/v1/verify/*`
