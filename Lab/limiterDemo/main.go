package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func main() {
	e := gin.Default()
	// 新建一个限速器，允许突发 10 个并发，限速 3rps，超过 500ms 就不再等待
	e.Use(NewLimiter(3, 10, 500*time.Millisecond))
	e.GET("ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	e.Run(":8080")
}
//主要就是使用了 sync.map  来为每一个 ip 创建一个 limiter，当然这个 key 也可以是其他的值，例如用户名等
func NewLimiter(r rate.Limit, b int, t time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		// 获取限速器
		// key 除了 ip 之外也可以是其他的，例如 header，user name 等
		key := c.ClientIP()
		l, _ := limiters.LoadOrStore(key, rate.NewLimiter(r, b))

		// 这里注意不要直接使用 gin 的 context 默认是没有超时时间的
		ctx, cancel := context.WithTimeout(c, t)
		defer cancel()

		if err := l.(*rate.Limiter).Wait(ctx); err != nil {
			// 这里先不处理日志了，如果返回错误就直接 429
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": err})
		}
		c.Next()
	}
