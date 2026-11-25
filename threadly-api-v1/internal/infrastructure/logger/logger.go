package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"threadly-api-v1/internal/config"
)

// Sampling matches zap.NewProductionConfig: per second, keep the first 100
// entries with the same level and message, then every 100th.
const (
	samplingTick       = time.Second
	samplingFirst      = 100
	samplingThereafter = 100
)

// New builds a logger that tees to every configured output. Stdout uses
// cfg.Format; the file output is always JSON so log shippers can parse it.
// The returned cleanup flushes buffered entries and closes the log file.
func New(environment string, cfg config.LogConfig) (*zap.Logger, func(), error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, nil, fmt.Errorf("parse log level: %w", err)
	}

	var (
		cores   []zapcore.Core
		closers []func() error
	)

	for _, output := range cfg.Outputs {
		switch output {
		case config.LogOutputStdout:
			cores = append(cores, zapcore.NewCore(newEncoder(cfg.Format), zapcore.Lock(os.Stdout), level))
		case config.LogOutputFile:
			fileWriter := &lumberjack.Logger{
				Filename:   cfg.File.Path,
				MaxSize:    cfg.File.MaxSizeMB,
				MaxBackups: cfg.File.MaxBackups,
				MaxAge:     cfg.File.MaxAgeDays,
				Compress:   cfg.File.Compress,
			}

			// lumberjack opens the file lazily; an empty write opens it now so
			// a bad path or missing permission fails startup instead of every log call.
			if _, err := fileWriter.Write(nil); err != nil {
				return nil, nil, fmt.Errorf("open log file %s: %w", cfg.File.Path, err)
			}

			cores = append(cores, zapcore.NewCore(newEncoder(config.LogFormatJSON), zapcore.AddSync(fileWriter), level))
			closers = append(closers, fileWriter.Close)
		default:
			return nil, nil, fmt.Errorf("unknown log output %q", output)
		}
	}

	core := zapcore.NewTee(cores...)
	options := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)}

	if environment == config.EnvironmentDevelopment {
		options = append(options, zap.Development())
	} else {
		core = zapcore.NewSamplerWithOptions(core, samplingTick, samplingFirst, samplingThereafter)
	}

	appLogger := zap.New(core, options...)

	cleanup := func() {
		// Sync on stdout returns EINVAL/ENOTTY on some platforms; nothing to act on at shutdown.
		_ = appLogger.Sync()

		for _, closeFile := range closers {
			if err := closeFile(); err != nil {
				fmt.Fprintf(os.Stderr, "close log file: %v\n", err)
			}
		}
	}

	return appLogger, cleanup, nil
}

func newEncoder(format string) zapcore.Encoder {
	if format == config.LogFormatConsole {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

// RequestLogger logs one line per HTTP request, replacing gin.Logger().
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		start := time.Now()

		context.Next()

		logger.Info("http request",
			zap.String("method", context.Request.Method),
			zap.String("path", context.Request.URL.Path),
			zap.Int("status", context.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", context.ClientIP()),
		)
	}
}
