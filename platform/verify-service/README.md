# Verify Service

职责：

- 接入 Tekton / Argo 事件
- 校验外部执行结果
- 把运行态 / 终态写回元数据层

当前状态：

- 已分出独立入口和基础回写接口
- `POST /api/v1/verify/argo/events`
- `POST /api/v1/verify/release/steps`
- `POST /api/v1/verify/tekton/events`
- `POST /api/v1/verify/tekton/steps`
- 当前可回写 Job / Manifest / Intent 的基础状态
- 可选共享鉴权：设置 `VERIFY_SERVICE_SHARED_TOKEN` 后，所有写接口都要求 `X-Devflow-Verify-Token`

最小回写语义：

- `argo/events`：更新 `Job.status`，并同步对应 release intent 状态
- `release/steps`：更新 `Job.steps` 中单个 step 的状态、进度、时间戳；缺失 step 时自动创建
- `tekton/events`：更新 `Manifest.status` / `pipeline_id`，并同步对应 build intent 状态
- `tekton/steps`：更新 `Manifest.steps` 中单个 task 的状态、taskRun、时间戳

这条 `release/steps` 接口用于承接 `devflow-controller` 一类观察器服务的状态同步，避免 controller 直接写 Mongo。
同一次 step 更新后，服务会自动从 steps 推导并更新 `job.status`。

Controller 对接说明见：

- `platform/verify-service/controller-integration.md`

示例：

```json
{
  "job_id": "67f000000000000000000001",
  "status": "Succeeded",
  "external_ref": "argocd/app/demo"
}
```

```json
{
  "job_id": "67f000000000000000000001",
  "step_name": "apply manifests",
  "status": "Running",
  "progress": 30,
  "message": "Application syncing"
}
```

```json
{
  "manifest_id": "67f000000000000000000002",
  "pipeline_id": "devflow-ci-run-abcde",
  "status": "Running"
}
```

建议端口：

- `VERIFY_SERVICE_PORT`
- `VERIFY_SERVICE_METRICS_PORT`
- `VERIFY_SERVICE_PPROF_PORT`

可选环境变量：

- `VERIFY_SERVICE_SHARED_TOKEN`

可观测性：

- 所有 verify 写接口都应产出 server span
- 对 Mongo / release-service / 外部回调方的调用都应产出 client span
- logs 必须带 `trace_id` / `span_id` 以及 `job_id` / `manifest_id` / `intent_id`
- profiling 建议优先覆盖 verify-service，因为它是状态汇聚点
- 上报的 OTel `service.name` 为 `verify-service`
- 规范见 `agents/reference/observability.md`
