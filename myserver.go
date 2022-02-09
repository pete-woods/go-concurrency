package go_concurrency

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/circleci/ex/httpserver/ginrouter"
	"github.com/circleci/ex/o11y"
	"github.com/gin-gonic/gin"
)

type MyServer struct {
	router *gin.Engine

	// mutable state managed by mutex
	mu      sync.Mutex
	counter int
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
	return s.counter
}

func (s *MyServer) getMyPage(c *gin.Context) {
	ctx := c.Request.Context()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	count := s.counter

	time.Sleep(10 * time.Millisecond)
	o11y.AddField(ctx, "count", count)
	c.String(http.StatusOK, fmt.Sprintf("%d", count))
}
