# Choosing a Go Web Framework

Go's standard library `net/http` is capable on its own, but most teams reach for a framework for routing, middleware, and ergonomics. Here's a comparison of the frameworks supported by Go Initializer.

## Gin

**GitHub Stars:** ~80k · **License:** MIT

The most widely used Go web framework. Fast router, minimal allocations, a large middleware ecosystem.

```go
r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
})
r.Run()
```

**Best for:** Teams that want battle-tested stability and a huge community to draw on.

---

## Echo

**GitHub Stars:** ~30k · **License:** MIT

Minimalist and highly performant. Clean API, first-class middleware support, built-in data binding and validation.

```go
e := echo.New()
e.GET("/", func(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
})
e.Logger.Fatal(e.Start(":1323"))
```

**Best for:** Teams that prefer a clean, opinionated API with strong documentation.

---

## Fiber

**GitHub Stars:** ~35k · **License:** MIT

Inspired by Express.js. Built on top of `fasthttp` for extreme throughput. Near-zero memory allocations.

```go
app := fiber.New()
app.Get("/", func(c *fiber.Ctx) error {
    return c.SendString("Hello, World!")
})
app.Listen(":3000")
```

**Best for:** High-throughput services where raw performance is critical.

---

## Chi

**GitHub Stars:** ~18k · **License:** MIT

`net/http`-compatible router with composable middleware. No external dependencies.

```go
r := chi.NewRouter()
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("welcome"))
})
http.ListenAndServe(":3000", r)
```

**Best for:** Teams that want to stay close to the stdlib and avoid framework lock-in.

---

## Summary

| Framework | Speed | Ecosystem | Stdlib-compatible |
|-----------|-------|-----------|-------------------|
| Gin       | ★★★★  | ★★★★★     | No                |
| Echo      | ★★★★  | ★★★★      | No                |
| Fiber     | ★★★★★ | ★★★       | No (fasthttp)     |
| Chi       | ★★★★  | ★★★       | Yes               |
