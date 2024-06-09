package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	// 2024/05/23 15:07:00 INFO Info message
	slog.Info("Info message")

	// {"time":"2024-05-23T15:08:06.349378+02:00","level":"INFO","msg":"Info message"}
	// {"time":"2024-05-23T15:08:06.349393+02:00","level":"WARN","msg":"Warning message"}
	// {"time":"2024-05-23T15:08:06.349398+02:00","level":"ERROR","msg":"Error message"}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// time=2024-05-23T15:09:16.355+02:00 level=INFO msg="Info message"
	// time=2024-05-23T15:09:16.355+02:00 level=WARN msg="Warning message"
	// time=2024-05-23T15:09:16.355+02:00 level=ERROR msg="Error message"
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// setting default global logger
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// {"time":"2024-05-23T15:10:25.353663+02:00","level":"INFO","msg":"Info message"}
	slog.Info("Info message")

	// context into logging

	// {"time":"2024-05-23T15:13:12.731011+02:00","level":"INFO","msg":"incoming request","method":"GET","time_taken_ms":158,"path":"/hello/world?q=search","status":200,"user_agent":"Googlebot/2.1 (+http://www.google.com/bot.html)"}
	logger.Info(
		"incoming request",
		"method", "GET",
		"time_taken_ms", 158,
		"path", "/hello/world?q=search",
		"status", 200,
		"user_agent", "Googlebot/2.1 (+http://www.google.com/bot.html)",
	)

	// {"time":"2024-05-23T15:15:20.285803+02:00","level":"INFO","msg":"incoming request","method":"GET","time_taken_ms":158,"path":"/hello/world?q=search","status":200,"user_agent":"Googlebot/2.1 (+http://www.google.com/bot.html)"}
	logger.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
		),
	)

	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
		),
	)

	// child loggers
	_ = logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", "1.21"),
		),
	)

	// change default logging level
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
	// {"time":"2024-05-23T15:24:13.285718+02:00","level":"DEBUG","msg":"Debug message"}
	logger.Debug("Debug message")

	// error logging
	ctx := context.Background()
	err := errors.New("something happened")
	// {"time":"2024-05-23T15:40:48.591876+02:00","level":"ERROR","msg":"upload failed","error":"something happened"}
	logger.ErrorContext(ctx, "upload failed", slog.Any("error", err))

	// sensitive values
	u := &User{
		ID:        "user-12234",
		FirstName: "Jan",
		LastName:  "Doe",
		Email:     "jan@example.com",
		Password:  "pass-12334",
	}

	// {"time":"2024-05-23T15:42:43.47524+02:00","level":"INFO","msg":"info","user":{"id":"user-12234","name":"Jan Doe"}}
	logger.Info("info", "user", u)

	// using 3rd party logging backends
	// One of Slog's major design goals is to provide a unified logging frontend (slog.Logger)
	// for Go applications while the backend (slog.Handler) remains customizable from program to program.
	zapL := zap.Must(zap.NewProduction())
	defer zapL.Sync()

	logger = slog.New(zapslog.NewHandler(zapL.Core(), nil))
	// {"level":"info","ts":1716472061.7753081,"msg":"incoming request","method":"GET","path":"/api/user","status":200}
	logger.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.String("path", "/api/user"),
		slog.Int("status", 200),
	)
	logger.Error("upload failed", slog.Any("error", err))
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// implement the `LogValuer` interface on the User struct
func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("name", u.FirstName+" "+u.LastName),
	)
}
