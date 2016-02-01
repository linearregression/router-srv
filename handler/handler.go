package handler

import (
	proto "github.com/micro/router-srv/proto/router"

	"golang.org/x/net/context"
)

type Router struct{}

func (r *Router) Stats(ctx context.Context, req *proto.StatsRequest, rsp *proto.StatsResponse) error {
	return nil
}

func (r *Router) Select(ctx context.Context, req *proto.SelectRequest, rsp *proto.SelectResponse) error {
	return nil
}

func (r *Router) Mark(ctx context.Context, req *proto.MarkRequest, rsp *proto.MarkResponse) error {
	return nil
}
