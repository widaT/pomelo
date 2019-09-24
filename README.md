# pomelo

轻量化go api框架


配置accesslog

```go
	s := pomelo.Default()
	s.Add("/", hello)
	s.Run()
```