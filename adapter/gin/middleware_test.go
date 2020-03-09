package gin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/bmizerany/assert"
	"github.com/gin-gonic/gin"
)

func initSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.FlowRule{
		{
			Resource:        "GET_/ping",
			MetricType:      flow.QPS,
			Count:           1,
			ControlBehavior: flow.Reject,
		},
		{
			Resource:        "/ping",
			MetricType:      flow.QPS,
			Count:           0,
			ControlBehavior: flow.Reject,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
		return
	}
}

func TestSentinelMiddleware(t *testing.T) {
	type args struct {
		opts []Option
		method string
		path string
		handler func(ctx *gin.Context)
		body io.Reader
	}
	type want struct {
		code int
	}
	var (
		tests = []struct {
			name string
			args args
			want want
		}{
			{
				name: "default get",
				args: args{
					opts: []Option{},
					method: http.MethodGet,
					path: "/ping",
					handler: func(ctx *gin.Context) {
						ctx.String(http.StatusOK, "ping")
					},
					body: nil,
				},
				want: want{
					code:http.StatusOK,
				},
			},
			{
				name: "customize resource extract",
				args: args{
					opts: []Option{
						WithResourceExtract(func(ctx *gin.Context) string {
							return ctx.Request.URL.Path
						}),
					},
					method: http.MethodGet,
					path: "/ping",
					handler: func(ctx *gin.Context) {
						ctx.String(http.StatusOK, "ping")
					},
					body: nil,
				},
				want: want{
					code:http.StatusTooManyRequests,
				},
			},
			{
				name: "customize block fallback",
				args: args{
					opts: []Option{
						WithResourceExtract(func(ctx *gin.Context) string {
							return ctx.Request.URL.Path
						}),
						WithBlockFallback(func(ctx *gin.Context) {
							ctx.String(http.StatusBadRequest, "block")
						}),
					},
					method: http.MethodGet,
					path: "/ping",
					handler: func(ctx *gin.Context) {
						ctx.String(http.StatusOK, "ping")
					},
					body: nil,
				},
				want: want{
					code:http.StatusBadRequest,
				},
			},
		}
	)
	initSentinel(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(SentinelMiddleware(tt.args.opts...))
			router.Handle(tt.args.method, tt.args.path, tt.args.handler)
			r := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}