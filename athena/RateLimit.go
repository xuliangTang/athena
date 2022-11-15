package athena

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
	"time"
)

// LimitOpt 限流配置
type LimitOpt struct {
	Interval int64 `mapstructure:"interval"` // 多长时间添加一次令牌
	Capacity int64 `mapstructure:"capacity"` // 令牌桶的容量
	Quantum  int64 `mapstructure:"quantum"`  // 到达定时器指定的时间，往桶里面加多少令牌
}

// LimitConfRules 限流规则map
type LimitConfRules struct {
	Rules map[string]*LimitOpt `mapstructure:"rateLimitRules"`
}

// ILimiter 限流方法实现定义
type ILimiter interface {
	Key(c *gin.Context) string
	GetBucket(key string) (*ratelimit.Bucket, bool)
	AddBucketsByUri(uri string, fillInterval, capacity, quantum int64) ILimiter
	AddBucketByConf() ILimiter
}

// UriLimiter 基于路由的限流
type UriLimiter struct {
	limiterBuckets map[string]*ratelimit.Bucket
	Rule           *LimitConfRules
}

func NewUriLimiter() ILimiter {
	return &UriLimiter{
		limiterBuckets: make(map[string]*ratelimit.Bucket),
	}
}

// Key 获取路由key
func (l *UriLimiter) Key(c *gin.Context) string {
	uri := c.Request.RequestURI
	index := strings.Index(uri, "?")
	if index == -1 {
		return uri
	}
	return uri[:index]
}

// GetBucket 根据路由key获取对应的buckets
func (l *UriLimiter) GetBucket(key string) (*ratelimit.Bucket, bool) {
	bucket, ok := l.limiterBuckets[key]
	return bucket, ok
}

// AddBucketsByUri 添加一条路由限流规则
func (l *UriLimiter) AddBucketsByUri(uri string, fillInterval, capacity, quantum int64) ILimiter {
	bucket := ratelimit.NewBucketWithQuantum(time.Second*time.Duration(fillInterval), capacity, quantum)
	l.limiterBuckets[uri] = bucket
	return l
}

// 读取相关配置
func (l *UriLimiter) getConf() *LimitConfRules {
	once := sync.Once{}
	rule := &LimitConfRules{}
	once.Do(func() {
		vp := viper.New()
		vp.SetConfigFile(FrameConf.AppPath + "/tests/application.yml")
		err := vp.ReadInConfig()
		if err != nil {
			log.Fatalln("read config.yaml error :", err)
		}
		errRule := vp.Unmarshal(&rule)
		if errRule != nil {
			log.Fatalln("unmarshal err :", errRule)
		}

		// 监控配置文件变化
		vp.WatchConfig()
		vp.OnConfigChange(func(in fsnotify.Event) {
			GlobalLimiter = NewUriLimiter().AddBucketByConf()
		})
	})
	return rule
}

// AddBucketByConf 通过配置文件批量添加限流规则
func (l *UriLimiter) AddBucketByConf() ILimiter {
	rule := l.getConf()
	for k, v := range rule.Rules {
		l.AddBucketsByUri(k, v.Interval, v.Capacity, v.Quantum)
	}
	return l
}

// GlobalLimiter 全局限流规则
var GlobalLimiter ILimiter

func init() {
	GlobalLimiter = NewUriLimiter().AddBucketByConf()
}

// RateLimit 限流中间件
type RateLimit struct {
}

func NewRateLimit() *RateLimit {
	return &RateLimit{}
}

func (this *RateLimit) OnRequest(ctx *gin.Context) error {
	key := GlobalLimiter.Key(ctx)
	if bucket, ok := GlobalLimiter.GetBucket(key); ok {
		count := bucket.TakeAvailable(1)
		if count == 0 {
			panic("rate limit")
			// ctx.AbortWithStatus(http.StatusForbidden)
			// return fmt.Errorf("rate limit")
		}
	}

	return nil
}
