package athena

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"time"
)

var logger *zap.Logger
var loggerOnce sync.Once

type LevelEnablerFunc func(lvl *zapcore.Level) bool

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

// Logger 获取日志对象
func Logger() *zap.Logger {
	loggerOnce.Do(func() {
		// 设置多log文件和轮转
		tops := getTops()
		var cores []zapcore.Core
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
		}

		for _, top := range tops {
			top := top

			lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return top.Lef(&lvl)
			})

			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   top.Filename,
				MaxSize:    top.Ropt.MaxSize,
				MaxBackups: top.Ropt.MaxBackups,
				MaxAge:     top.Ropt.MaxAge,
				Compress:   top.Ropt.Compress,
			})

			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg.EncoderConfig),
				zapcore.AddSync(w),
				lv,
			)
			cores = append(cores, core)
		}

		logger = zap.New(zapcore.NewTee(cores...))
		defer logger.Sync() // flushes buffer, if any
	})

	return logger
}

func getTops() []TeeOption {
	var tops = []TeeOption{
		{
			Filename: FrameConf.AppPath + "/storage/logs/access.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl <= zapcore.InfoLevel
			},
		},
		{
			Filename: FrameConf.AppPath + "/storage/logs/error.log",
			Ropt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl > zapcore.InfoLevel
			},
		},
	}

	return tops
}
