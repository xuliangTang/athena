package athena

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"sync"
	"time"
)

var logging *LoggingImpl
var logger *zap.Logger
var loggerOnce sync.Once

type LoggingImpl struct {
	*zap.Logger
}

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
func Logger() *LoggingImpl {
	loggerOnce.Do(func() {
		// 设置多log文件和轮转
		tops := getTops()
		var cores []zapcore.Core
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}

		for _, top := range tops {
			lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return top.Lef(&lvl)
			})

			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   top.Filename,
				MaxSize:    top.Ropt.MaxSize,
				MaxBackups: top.Ropt.MaxBackups,
				MaxAge:     top.Ropt.MaxAge,
				Compress:   top.Ropt.Compress,
				LocalTime:  true,
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
		logging = &LoggingImpl{logger}
	})

	return logging
}

// 重写父类方法
func (this *LoggingImpl) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("stack", this.GetStack()))

	this.Logger.Error(msg, fields...)
}

// GetStack 获取堆栈信息
func (this *LoggingImpl) GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("==> %s\n", string(buf[:n]))
}

func getTops() []TeeOption {
	var tops = []TeeOption{
		{
			Filename: FrameConf.AppPath + FrameConf.LogAccess.FilePath,
			Ropt: RotateOptions{
				MaxSize:    FrameConf.LogAccess.MaxSize,    // 日志大小限制，单位MB
				MaxAge:     FrameConf.LogAccess.MaxAge,     // 历史日志文件保留天数
				MaxBackups: FrameConf.LogAccess.MaxBackups, // 最大保留历史日志数量
				Compress:   false,                          // 历史日志文件压缩标识
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl <= zapcore.InfoLevel
			},
		},
		{
			Filename: FrameConf.AppPath + FrameConf.LogError.FilePath,
			Ropt: RotateOptions{
				MaxSize:    FrameConf.LogError.MaxSize,
				MaxAge:     FrameConf.LogError.MaxAge,
				MaxBackups: FrameConf.LogError.MaxBackups,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl > zapcore.InfoLevel
			},
		},
	}

	return tops
}
