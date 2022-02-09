package go_concurrency

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/circleci/ex/httpserver/ginrouter"
	"github.com/circleci/ex/o11y"
	"github.com/gin-gonic/gin"
)

type MyServer struct {
	router  *gin.Engine
	counter int
}

func NewMyServer(ctx context.Context) *MyServer {
	r := ginrouter.Default(ctx, "my-server")
	s := &MyServer{router: r}
	time.Sleep(10 * time.Millisecond) // simulate database access
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
	s.counter++
	o11y.AddField(ctx, "count", s.counter)
	c.String(http.StatusOK, fmt.Sprintf("%d", s.counter))
}
