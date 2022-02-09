package go_concurrency

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/circleci/ex/httpserver/ginrouter"
	"github.com/circleci/ex/o11y"
	"github.com/gin-gonic/gin"
)

type MyServer struct {
	router *gin.Engine

	// mutable state
	counter int64
}

func NewMyServer(ctx context.Context) *MyServer {
	r := ginrouter.Default(ctx, "my-server")
	s := &MyServer{router: r}
	r.GET("/mypage", s.getMyPage)
	return s
}

func (s *MyServer) Handler() http.Handler {
	return s.router
}

func (s *MyServer) Counter() int {
	return int(s.counter)
}

func (s *MyServer) getMyPage(c *gin.Context) {
	ctx := c.Request.Context()
	count := atomic.AddInt64(&s.counter, 1)

	time.Sleep(10 * time.Millisecond)
	o11y.AddField(ctx, "count", count)
	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}
