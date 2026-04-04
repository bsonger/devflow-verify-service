# 接口规范

`devflow-verify-service` 只暴露 `/api/v1/verify/*`。

## 接口列表

- `GET /api/v1/verify/healthz`
  - 用途：verify 服务健康检查
  - 返回：`200`
- `POST /api/v1/verify/argo/events`
  - 用途：回写 `Job` 级发布状态
  - 关键字段：`job_id`、`status`、`intent_id`、`external_ref`
- `POST /api/v1/verify/tekton/events`
  - 用途：回写 `Manifest` 级构建状态
  - 关键字段：`manifest_id`、`status`、`pipeline_id`、`intent_id`、`external_ref`
- `POST /api/v1/verify/tekton/steps`
  - 用途：回写 `Manifest.steps`
  - 关键字段：`manifest_id`、`pipeline_id`、`task_name`、`task_run`、`status`
- `POST /api/v1/verify/release/steps`
  - 用途：回写 `Job.steps`
  - 关键字段：`job_id`、`step_name`、`status`、`progress`

## 认证

- `/api/v1/verify/*` 写接口使用 `X-Devflow-Verify-Token`
- 若未设置 `VERIFY_SERVICE_SHARED_TOKEN`，本地环境可无 token 访问
- 若设置了 `VERIFY_SERVICE_SHARED_TOKEN`，写接口必须校验通过才允许写入

## 错误语义

- `400`：请求体缺失、资源 ID 非法、必要字段缺失
- `401`：共享 token 校验失败
- `500`：Mongo 更新失败或内部写回异常

## 非目标

- 不提供分页接口
- 不提供 `Project`、`Application`、`Configuration`、`Manifest`、`Job`、`Intent` 的对外 CRUD
