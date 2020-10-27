# echo_health

an Echo middleware for health checking

## Usage

```go
e := echo.New()
e.Use(echo_heath.New(echo_health.ReadinessCheckFunc(func() error {
  return DB.Ping()
})))
```

```
GET /health/alive

GET /health/ready
```

## Credits

Guo Y.K., MIT License
