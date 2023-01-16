# Logecho

A package to use as logger pkg and middleware to [echo](https://echo.labstack.com/) framework. It wraps [zap logger](https://pkg.go.dev/go.uber.org/zap) into a logecho package

# Install

```shell
go get github.com/jeanmolossi/logecho
```

# Usage example

Basic usage

```go
package main

import (
	"github.com/jeanmolossi/logecho"
	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    // apply logecho middleware
    e.Use(logecho.Middleware())
    // ...
}
```

Templated usage

```go
package main

import (
	"github.com/jeanmolossi/logecho"
	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    // apply logecho middleware
    e.Use(
        logecho.MiddlewareWithTemplate(logecho.Fields{
            "host":           logecho.Host,
            "request.method": logecho.Method,
            "request.path":   logecho.Path,
            "latency":        logecho.LatencyString,
            "request-id":     logecho.RequestID,
            // x-user-id is the wanted header and user-id is 
            // fallback when x-user-id is not present
            "user-id":        logecho.Header("x-user-id", "user-id"),
        }),
    )
    // ...
}
```

## Logging inside a handler

```go
package handlers

import (
	"github.com/jeanmolossi/logecho"
	"github.com/labstack/echo/v4"
)

// recover logecho.Logger singleton
var log = logecho.Logger

func handler(c echo.Context) error {
    // log passing c (echo.Context) to write 
    // fields to that previous configured on middleware
    //
    // Output example:
    //  {"timestamp": "2012-01-31T00:00:00-0300","level": "info","message": "started handler request", "host": "localhost:8080","user-id": "1", ...}
    log.Info(c, "started handler request")

    // ... your logic

    // another log with same fields previous configured
    log.Info(c, "end handler request")

    return c.JSON(http.StatusOK, echo.Map{"message": "done"})
}
```

## Additional info

You can pass down only `echo.Context` to your sub calls inside your handle and use the same instance of `logecho.Logger`. It was designed to be thread-safe (that was the try)