# tracex 适配器

适配器模块保持 tracex 核心零第三方依赖，按需引入框架/库适配。

## webx（服务端中间件）

```go
import txwebx "github.com/lcylpzls/tracex/adapters/webx"

s.UseGlobalMiddleware(txwebx.Middleware(m))
```

## dbx（数据库查询）

```go
import txdbx "github.com/lcylpzls/tracex/adapters/dbx"

db, _ := dbx.Open(ctx, "mysql", dsn, dbx.WithTraceHook(txdbx.NewHook(m)))
```

Exec / Query / QueryRow 自动生成 dbx.* span。

## jobx（任务执行）

```go
import txjobx "github.com/lcylpzls/tracex/adapters/jobx"

d, _ := jobx.NewDispatcher(jobx.WithTraceHook(txjobx.NewHook(m)))
```

任务处理器执行自动生成 jobx.execute span。

## cachex（回源加载）

```go
import txcachex "github.com/lcylpzls/tracex/adapters/cachex"

c, _ := cachex.New(cachex.WithTraceHook(txcachex.NewHook(m)))
```

GetOrSet 实际回源时自动生成 cachex.load span，命中缓存不埋点。

## httpx（出站请求）

```go
client, _ := httpx.New(httpx.WithRoundTripperWrapper(m.RoundTripper))
```

保留 httpx 协议选择（HTTP/1.1/2/3），自动携带 traceparent 并生成
客户端 span。
