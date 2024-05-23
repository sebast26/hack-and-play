package main

import (
	"errors"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	// {"level":"info","ts":1716395546.635746,"caller":"zap-logging/main.go:13","msg":"Hello from Zap logger!"}
	logger.Info("Hello from Zap logger!")

	// using global logger
	zap.L().Info("Hello from Zap!")

	// low-level Logger API
	// {"level":"info","ts":1716395546.6359098,"caller":"zap-logging/main.go:19","msg":"User logged in","username":"sebastian","userID":1234567,"provider":"google"}
	logger.Info("User logged in",
		zap.String("username", "sebastian"),
		zap.Int("userID", 1234567),
		zap.String("provider", "google"),
	)

	// higher-level Logger API (SugaredLogger)
	sugar := logger.Sugar()

	sugar.Info("Hello from Zap logger!")
	// {"level":"info","ts":1716411566.73898,"caller":"zap-logging/main.go:31","msg":"Hello from Zap logger"}
	sugar.Infoln("Hello", "from", "Zap", "logger")
	sugar.Infof(
		"Hello from Zap logger! The time is %s",
		time.Now().Format("03:04 AM"),
	)
	sugar.Infow("User logged in",
		"username", "johndoe",
		"userid", 123456,
		zap.String("provider", "google"),
	)
	sugar.Desugar() // in-expensive

	// creating a custom logger
	clogger := createLogger()
	defer clogger.Sync()

	// {"level":"info","timestamp":"2024-05-22T23:11:30.498+0200","caller":"zap-logging/main.go:54","msg":"Hello from Zap!","pid":3003}
	clogger.Info("Hello from Zap!")

	// adding context to logs

	childLogger := logger.With(
		zap.String("service", "userService"),
		zap.String("requestID", "abc123"),
	)

	// {"level":"info","ts":1716412777.4505112,"caller":"zap-logging/main.go:64","msg":"user registration successfil","service":"userService","requestID":"abc123","username":"john.doe","email":"john@example.com"}
	childLogger.Info("user registration successfil",
		zap.String("username", "john.doe"),
		zap.String("email", "john@example.com"),
	)

	// {"level":"info","ts":1716412777.4505222,"caller":"zap-logging/main.go:69","msg":"redirectign user to adming dashboard","service":"userService","requestID":"abc123"}
	childLogger.Info("redirectign user to adming dashboard")

	// logging errors
	// {"level":"error","ts":1716467396.108449,"caller":"zap-logging/main.go:75","msg":"failed to perform an operation","operation":"someOperation","error":"something happened","retryAttempts":3,"user":"john.doe","stacktrace":"main.main\n\t/Users/seba/projects/priv/code/hack-and-play/golang/zap-logging/main.go:75\nruntime.main\n\t/usr/local/Cellar/go/1.21.3/libexec/src/runtime/proc.go:267"}
	logger.Error("failed to perform an operation",
		zap.String("operation", "someOperation"),
		zap.Error(errors.New("something happened")), // the key will be error
		zap.Int("retryAttempts", 3),
		zap.String("user", "john.doe"),
	)

	// sensitive fields

}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SensitiveFieldEncoder struct {
	zapcore.Encoder
	cfg zapcore.EncoderConfig
}

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	_ = NewSensitiveFieldsEncoder(encoderCfg)

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

// EncodeEntry is called for every log line to be emitted so it needs to be
// as efficient as possible so that you don't negate the speed/memory advantages
// of Zap
func (e *SensitiveFieldEncoder) EncodeEntry(
	entry zapcore.Entry,
	fields []zapcore.Field,
) (*buffer.Buffer, error) {
	filtered := make([]zapcore.Field, 0, len(fields))

	for _, field := range fields {
		user, ok := field.Interface.(User)
		if ok {
			user.Email = "[REDACTED]"
			field.Interface = user
		}

		filtered = append(filtered, field)
	}

	return e.Encoder.EncodeEntry(entry, filtered)
}

func NewSensitiveFieldsEncoder(config zapcore.EncoderConfig) zapcore.Encoder {
	encoder := zapcore.NewJSONEncoder(config)
	return &SensitiveFieldEncoder{encoder, config}
}
