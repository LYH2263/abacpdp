# go-abacpdp

Attribute-Based Access Control **策略判决点（PDP）**：装载 PolicySet，绑定 Subject / Resource / Action / Environment，经属性匹配与 DenyOverrides 等合并算法输出 Permit/Deny 与 Obligations。

## 与其它仓的差异

- ≠ `go-apigatex`：无反向代理 / 路由改写
- ≠ `go-flagship`：非百分比功能开关
- ≠ `go-ticketgate`：无票据格式

## 构建与测试

```text
go build ./...
go test ./... -count=1
go run ./cmd/pdpd -addr :8112
```

管理页：编辑试属性，查看 Permit/Deny 与义务。

## 主类型

- `PDP` — `Load` / `Evaluate` / `EvaluateContext` / `Close`
- `AttrBag` — 属性袋（可 Clone）
- `Decision` / `Obligation` / `Effect`
