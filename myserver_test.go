package go_concurrency

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/circleci/ex/httpclient"
	"github.com/circleci/ex/o11y"
	"github.com/circleci/ex/testing/testcontext"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestMyServer(t *testing.T) {
	ctx := context.Background()
	s := NewMyServer(ctx)

	srv := httptest.NewServer(s.Handler())
	t.Cleanup(srv.Close)

	c := httpclient.New(httpclient.Config{
		Name:    "client",
		BaseURL: srv.URL,
	})

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for i := 0; i < 50; i++ {
				err := testCall(c, ctx)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	assert.Assert(t, g.Wait())

	o11y.Log(testcontext.Background(), "final count", o11y.Field("count", s.Counter()))
}

func testCall(c *httpclient.Client, ctx context.Context) (err error) {
	ctx, span := o11y.StartSpan(ctx, "testcall")
	defer o11y.End(span, &err)

	count := ""
	err = c.Call(ctx, httpclient.Request{
		Method:  "GET",
		Route:   "/mypage",
		Decoder: httpclient.NewStringDecoder(&count),
	})
	if err != nil {
		return err
	}
	span.AddField("count", count)
	return nil
}
