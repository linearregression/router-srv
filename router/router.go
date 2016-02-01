package router

import (
	"sync"

	proto "github.com/micro/go-platform/router/proto"

	"golang.org/x/net/context"
)

var (
	DefaultRouter = newRouter()
	DefaultWindow = 10

	StatsTopic = "platform.router.stats"
)

type router struct {
	sync.Mutex
	stats  map[string][]*proto.Stats
	window int
}

func newRouter() *router {
	return &router{
		stats:  make(map[string][]*proto.Stats),
		window: DefaultWindow,
	}
}

func (r *router) ProcessStats(ctx context.Context, stats *proto.Stats) error {
	r.Lock()
	defer r.Unlock()

	if stats.Service == nil || len(stats.Service.Nodes) == 0 {
		return nil
	}

	rs, ok := r.stats[stats.Service.Nodes[0].Id]
	if !ok {
		rs = []*proto.Stats{}
	}

	rs = append(rs, stats)
	if len(rs) > r.window {
		rs = rs[:r.window]
	}

	r.stats[stats.Service.Nodes[0].Id] = rs

	// Do some actual stats processing here.

	return nil
}

func ProcessStats(ctx context.Context, stats *proto.Stats) error {
	return DefaultRouter.ProcessStats(ctx, stats)
}
