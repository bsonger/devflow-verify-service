# API Contract

本文件是 agent 参考摘要，正式规范以 `docs/api-spec.md` 为准。

当前仓库只允许：

- `Verify`

禁止出现：

- `Project`
- `Application`
- `Configuration`
- `Manifest`
- `Release`
- `Intent`

要求：

- handler、router、Swagger 三者一致
- 所有 public path 必须位于 `/api/v1/verify/*`
- 鉴权规则和错误语义更新到 `docs/api-spec.md`
