package lib

import (
	"bytes"
	"github.com/xuliangTang/athena/athena/config"
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
		trees := getZapTree()
		var cores []zapcore.Core
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}

		for i, _ := range trees {
			topIndex := i
			lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return trees[topIndex].Lef(&lvl)
			})

			ws := zapcore.AddSync(&lumberjack.Logger{
				Filename:   trees[i].Filename,
				MaxSize:    trees[i].Ropt.MaxSize,
				MaxBackups: trees[i].Ropt.MaxBackups,
				MaxAge:     trees[i].Ropt.MaxAge,
				Compress:   trees[i].Ropt.Compress,
				LocalTime:  true,
			})

			core := zapcore.NewCore(
				zapcore.NewJSONEncoder(cfg.EncoderConfig),
				ws,
				lv,
			)
			cores = append(cores, core)
		}

		logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller())
		defer logger.Sync()
		logging = &LoggingImpl{logger}
	})

	return logging
}

// 重写父类方法
func (this *LoggingImpl) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("stack", this.GetStack(10)))

	this.Logger.Error(msg, fields...)
}

// GetStack 获取堆栈信息
func (this *LoggingImpl) GetStack(kb int) string {
	/*var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("==> %s\n", string(buf[:n]))*/

	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) // 4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}

// 所有日志类型
func getZapTree() []TeeOption {
	var tops = []TeeOption{
		{
			Filename: config.AppConf.AppPath + config.AppConf.Logging.LogAccess.FilePath, // 输出文件目录
			Ropt: RotateOptions{
				MaxSize:    config.AppConf.Logging.LogAccess.MaxSize,    // 日志大小限制，单位MB
				MaxAge:     config.AppConf.Logging.LogAccess.MaxAge,     // 历史日志文件保留天数
				MaxBackups: config.AppConf.Logging.LogAccess.MaxBackups, // 最大保留历史日志数量
				Compress:   false,                                       // 历史日志文件压缩标识
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl <= zapcore.InfoLevel
			},
		},
		{
			Filename: config.AppConf.AppPath + config.AppConf.Logging.LogError.FilePath,
			Ropt: RotateOptions{
				MaxSize:    config.AppConf.Logging.LogError.MaxSize,
				MaxAge:     config.AppConf.Logging.LogError.MaxAge,
				MaxBackups: config.AppConf.Logging.LogError.MaxBackups,
				Compress:   false,
			},
			Lef: func(lvl *zapcore.Level) bool {
				return *lvl > zapcore.InfoLevel
			},
		},
	}

	return tops
}
