package router

import (
	"errors"
	"sync"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	proto "github.com/micro/go-platform/router/proto"
	"github.com/micro/router-srv/label"
	"golang.org/x/net/context"
)

var (
	DefaultRouter = newRouter()
	DefaultWindow = 10
	DefaultExpiry = 60 // seconds

	StatsTopic = "platform.router.stats"
)

type router struct {
	sync.RWMutex
	stats  map[string][]*proto.Stats
	window int

	r registry.Registry

	mtx      sync.Mutex
	pointers map[string]int
}

func newRouter() *router {
	return &router{
		stats:    make(map[string][]*proto.Stats),
		window:   DefaultWindow,
		r:        registry.DefaultRegistry,
		pointers: make(map[string]int),
	}
}

func (r *router) Init(s micro.Service) {
	r.r = s.Server().Options().Registry
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

func (r *router) Select(service string) ([]*proto.Service, error) {
	if len(service) == 0 {
		return nil, errors.New("invalid service")
	}

	// TODO: cache records or watch for updates
	services, err := r.r.GetService(service)
	if err != nil {
		return nil, err
	}

	srvLen := len(services)

	if srvLen == 0 {
		return nil, selector.ErrNotFound
	}

	// grab dynamic runtime labels
	// argh all the labels
	// bad bad bad, don't use 1000
	// TODO: fix
	labels, err := label.Search(service, "", 1000, 0)
	if err != nil {
		return nil, err
	}

	// TODO: use stats to assign weights to nodes
	// rather than just arbitrary pointer selection

	// get the pointer and increment it
	r.mtx.Lock()
	pointer := r.pointers[service]
	pointer++
	r.pointers[service] = pointer
	r.mtx.Unlock()

	servs := make([]*proto.Service, srvLen)

	// create starting point based on pointer and length of services
	i := pointer % srvLen

	// iterate through the length of the services
	for j := 0; j < srvLen; j++ {
		// if the pointer has dropped past length rotate back through
		if i >= srvLen {
			i = 0
		}

		// apply the label
		label.Apply(labels, services[i])

		// save the service, pass pointer toProto so it can rotate too
		servs[j] = toProto(services[i], pointer)
		i++
	}

	return servs, nil
}

func (r *router) Stats(service, nodeId string) ([]*proto.Stats, error) {
	if len(service) == 0 {
		return nil, errors.New("invalid service")
	}

	var stats []*proto.Stats

	r.RLock()
	defer r.RUnlock()

	for n, stat := range r.stats {
		// any stats
		if len(stat) == 0 {
			continue
		}

		// do we have node id?
		if len(nodeId) > 0 && n == nodeId {
			// grab the last stat
			stats = append(stats, stat[len(stat)-1])
			return stats, nil
		}

		if st := stat[len(stat)-1]; service == st.Service.Name {
			stats = append(stats, st)
		}
	}

	return stats, nil
}

func Init(s micro.Service) {
	DefaultRouter.Init(s)
}

func Select(service string) ([]*proto.Service, error) {
	return DefaultRouter.Select(service)
}

func ProcessStats(ctx context.Context, stats *proto.Stats) error {
	return DefaultRouter.ProcessStats(ctx, stats)
}

func Stats(service, nodeId string) ([]*proto.Stats, error) {
	return DefaultRouter.Stats(service, nodeId)
}
