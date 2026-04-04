# Verify Service Notes

## 当前定位

- 对外只提供 `Verify` 写回接口
- 内部只保留写回 `Manifest`、`Job`、`Intent` 所需的最小能力

## 已迁移内容

- 启动骨架迁移到 `../devflow-service-common/bootstrap`
- 观测 server 和公共 Gin 中间件迁移到 `../devflow-service-common`

## 参考

- 仓库级文档优先看 `../README.md` 与 `../docs/`
