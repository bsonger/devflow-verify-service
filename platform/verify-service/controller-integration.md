# Controller Integration

这份说明用于把 `devflow-controller` 一类观察器从“直接写 Mongo”迁到“调用 verify-service”。

## 目标边界

- controller 只监听外部系统状态
- controller 不直接持有 metadata Mongo 写权限
- controller 统一调用 `verify-service`

## 鉴权

如果设置了 `VERIFY_SERVICE_SHARED_TOKEN`：

- 请求头必须带 `X-Devflow-Verify-Token`
- `healthz` 不需要鉴权

## 事件映射

### 1. Argo Application 状态

用途：

- 更新 `job.status`
- 更新 apply 类 step，例如 `apply manifests`

接口：

- `POST /api/v1/verify/argo/events`
- `POST /api/v1/verify/release/steps`

典型 payload：

```json
{
  "job_id": "67f000000000000000000001",
  "status": "Running",
  "external_ref": "argocd/application/demo",
  "message": "Application syncing"
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

### 2. Deployment 状态

用途：

- 更新普通发布下的 deploy 类 step，例如 `deploy ready`

接口：

- `POST /api/v1/verify/release/steps`

典型 payload：

```json
{
  "job_id": "67f000000000000000000001",
  "step_name": "deploy ready",
  "status": "Running",
  "progress": 60,
  "message": "pods ready 3/5"
}
```

### 3. Canary Rollout 状态

用途：

- 更新金丝雀流量阶段 step，例如 `canary 10% traffic`

接口：

- `POST /api/v1/verify/release/steps`

典型 payload：

```json
{
  "job_id": "67f000000000000000000001",
  "step_name": "canary 30% traffic",
  "status": "Running",
  "progress": 50,
  "message": "canary weight 20% (target 30%)"
}
```

### 4. Blue-Green Rollout 状态

用途：

- 更新 `green ready`
- 更新 `switch traffic`

接口：

- `POST /api/v1/verify/release/steps`

典型 payload：

```json
{
  "job_id": "67f000000000000000000001",
  "step_name": "green ready",
  "status": "Succeeded",
  "progress": 100,
  "message": "green pods ready 5/5"
}
```

```json
{
  "job_id": "67f000000000000000000001",
  "step_name": "switch traffic",
  "status": "Succeeded",
  "progress": 100,
  "message": "traffic switched to green"
}
```

## 服务侧语义

- `release-service` 创建 Job 时会自动初始化默认 step 模板
- `verify-service` 更新 step 时，如果 step 不存在，会自动补创建
- `verify-service` 会根据所有 steps 自动收敛 `job.status`
- `argo/events` 仍可直接写 Job 级状态，作为更直接的外部事实来源
