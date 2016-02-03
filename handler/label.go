package handler

import (
	"github.com/micro/go-micro/errors"
	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/label"
	proto "github.com/micro/router-srv/proto/label"

	"golang.org/x/net/context"
)

type Label struct{}

func (r *Label) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Read", "invalid id")
	}

	l, err := label.Read(req.Id)
	if err != nil {
		if err == db.ErrNotFound {
			return errors.NotFound("go.micro.srv.router.Read", err.Error())
		}
		return errors.InternalServerError("go.micro.srv.router.Read", err.Error())
	}

	rsp.Label = l

	return nil
}

func (r *Label) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if req.Label == nil {
		return errors.BadRequest("go.micro.srv.router.Create", "invalid label")
	}

	if len(req.Label.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Create", "invalid id")
	}

	if len(req.Label.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Create", "invalid service")
	}

	if len(req.Label.Key) == 0 {
		return errors.BadRequest("go.micro.srv.router.Create", "invalid key")
	}

	if req.Label.Weight < 0 || req.Label.Weight > 100 {
		return errors.BadRequest("go.micro.srv.router.Create", "invalid weight, must be 0 to 100")
	}

	if err := label.Create(req.Label); err != nil {
		return errors.InternalServerError("go.micro.srv.router.Create", err.Error())
	}

	return nil
}

func (r *Label) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	if req.Label == nil {
		return errors.BadRequest("go.micro.srv.router.Update", "invalid label")
	}

	if len(req.Label.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Update", "invalid id")
	}

	if len(req.Label.Service) == 0 {
		return errors.BadRequest("go.micro.srv.router.Update", "invalid service")
	}

	if len(req.Label.Key) == 0 {
		return errors.BadRequest("go.micro.srv.router.Update", "invalid key")
	}

	if req.Label.Weight < 0 || req.Label.Weight > 100 {
		return errors.BadRequest("go.micro.srv.router.Update", "invalid weight, must be 0 to 100")
	}

	if err := label.Update(req.Label); err != nil {
		return errors.InternalServerError("go.micro.srv.router.Update", err.Error())
	}

	return nil
}

func (r *Label) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.router.Delete", "invalid id")
	}

	if err := label.Delete(req.Id); err != nil {
		if err == db.ErrNotFound {
			return nil
		}
		return errors.InternalServerError("go.micro.srv.router.Delete", err.Error())
	}

	return nil
}

func (r *Label) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Offset <= 0 {
		req.Offset = 0
	}

	labels, err := label.Search(req.Service, req.Key, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.router.Search", err.Error())
	}

	rsp.Labels = labels

	return nil
}
