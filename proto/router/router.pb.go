// Code generated by protoc-gen-go.
// source: github.com/micro/router-srv/proto/router/router.proto
// DO NOT EDIT!

/*
Package router is a generated protocol buffer package.

It is generated from these files:
	github.com/micro/router-srv/proto/router/router.proto

It has these top-level messages:
	Filter
	Expression
	StatsRequest
	StatsResponse
	SelectRequest
	SelectResponse
	MarkRequest
	MarkResponse
*/
package router

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import router1 "github.com/micro/go-platform/router/proto"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Filter struct {
	Version  *Expression   `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	Metadata []*Expression `protobuf:"bytes,2,rep,name=metadata" json:"metadata,omitempty"`
}

func (m *Filter) Reset()                    { *m = Filter{} }
func (m *Filter) String() string            { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()               {}
func (*Filter) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Filter) GetVersion() *Expression {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *Filter) GetMetadata() []*Expression {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type Expression struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	// equals, not equals
	Operator string `protobuf:"bytes,3,opt,name=operator" json:"operator,omitempty"`
}

func (m *Expression) Reset()                    { *m = Expression{} }
func (m *Expression) String() string            { return proto.CompactTextString(m) }
func (*Expression) ProtoMessage()               {}
func (*Expression) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StatsRequest struct {
	Service string `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	NodeId  string `protobuf:"bytes,2,opt,name=node_id" json:"node_id,omitempty"`
}

func (m *StatsRequest) Reset()                    { *m = StatsRequest{} }
func (m *StatsRequest) String() string            { return proto.CompactTextString(m) }
func (*StatsRequest) ProtoMessage()               {}
func (*StatsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type StatsResponse struct {
	Stats []*router1.Stats `protobuf:"bytes,1,rep,name=stats" json:"stats,omitempty"`
}

func (m *StatsResponse) Reset()                    { *m = StatsResponse{} }
func (m *StatsResponse) String() string            { return proto.CompactTextString(m) }
func (*StatsResponse) ProtoMessage()               {}
func (*StatsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *StatsResponse) GetStats() []*router1.Stats {
	if m != nil {
		return m.Stats
	}
	return nil
}

type SelectRequest struct {
	Service string    `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Filter  []*Filter `protobuf:"bytes,2,rep,name=filter" json:"filter,omitempty"`
}

func (m *SelectRequest) Reset()                    { *m = SelectRequest{} }
func (m *SelectRequest) String() string            { return proto.CompactTextString(m) }
func (*SelectRequest) ProtoMessage()               {}
func (*SelectRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SelectRequest) GetFilter() []*Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type SelectResponse struct {
	Nodes []*router1.Node `protobuf:"bytes,1,rep,name=nodes" json:"nodes,omitempty"`
}

func (m *SelectResponse) Reset()                    { *m = SelectResponse{} }
func (m *SelectResponse) String() string            { return proto.CompactTextString(m) }
func (*SelectResponse) ProtoMessage()               {}
func (*SelectResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SelectResponse) GetNodes() []*router1.Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type MarkRequest struct {
	Service string        `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Node    *router1.Node `protobuf:"bytes,2,opt,name=node" json:"node,omitempty"`
	Error   string        `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
}

func (m *MarkRequest) Reset()                    { *m = MarkRequest{} }
func (m *MarkRequest) String() string            { return proto.CompactTextString(m) }
func (*MarkRequest) ProtoMessage()               {}
func (*MarkRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MarkRequest) GetNode() *router1.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

type MarkResponse struct {
}

func (m *MarkResponse) Reset()                    { *m = MarkResponse{} }
func (m *MarkResponse) String() string            { return proto.CompactTextString(m) }
func (*MarkResponse) ProtoMessage()               {}
func (*MarkResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*Filter)(nil), "Filter")
	proto.RegisterType((*Expression)(nil), "Expression")
	proto.RegisterType((*StatsRequest)(nil), "StatsRequest")
	proto.RegisterType((*StatsResponse)(nil), "StatsResponse")
	proto.RegisterType((*SelectRequest)(nil), "SelectRequest")
	proto.RegisterType((*SelectResponse)(nil), "SelectResponse")
	proto.RegisterType((*MarkRequest)(nil), "MarkRequest")
	proto.RegisterType((*MarkResponse)(nil), "MarkResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Router service

type RouterClient interface {
	Stats(ctx context.Context, in *StatsRequest, opts ...client.CallOption) (*StatsResponse, error)
	Select(ctx context.Context, in *SelectRequest, opts ...client.CallOption) (*SelectResponse, error)
	Mark(ctx context.Context, in *MarkRequest, opts ...client.CallOption) (*MarkResponse, error)
}

type routerClient struct {
	c           client.Client
	serviceName string
}

func NewRouterClient(serviceName string, c client.Client) RouterClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "router"
	}
	return &routerClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *routerClient) Stats(ctx context.Context, in *StatsRequest, opts ...client.CallOption) (*StatsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Router.Stats", in)
	out := new(StatsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerClient) Select(ctx context.Context, in *SelectRequest, opts ...client.CallOption) (*SelectResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Router.Select", in)
	out := new(SelectResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerClient) Mark(ctx context.Context, in *MarkRequest, opts ...client.CallOption) (*MarkResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Router.Mark", in)
	out := new(MarkResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Router service

type RouterHandler interface {
	Stats(context.Context, *StatsRequest, *StatsResponse) error
	Select(context.Context, *SelectRequest, *SelectResponse) error
	Mark(context.Context, *MarkRequest, *MarkResponse) error
}

func RegisterRouterHandler(s server.Server, hdlr RouterHandler) {
	s.Handle(s.NewHandler(&Router{hdlr}))
}

type Router struct {
	RouterHandler
}

func (h *Router) Stats(ctx context.Context, in *StatsRequest, out *StatsResponse) error {
	return h.RouterHandler.Stats(ctx, in, out)
}

func (h *Router) Select(ctx context.Context, in *SelectRequest, out *SelectResponse) error {
	return h.RouterHandler.Select(ctx, in, out)
}

func (h *Router) Mark(ctx context.Context, in *MarkRequest, out *MarkResponse) error {
	return h.RouterHandler.Mark(ctx, in, out)
}

var fileDescriptor0 = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x52, 0x61, 0x8b, 0xda, 0x40,
	0x10, 0xd5, 0x6a, 0xa2, 0x9d, 0x18, 0x2d, 0xdb, 0x96, 0x4a, 0x68, 0x41, 0x16, 0x5a, 0x84, 0xe2,
	0x5a, 0x2c, 0x2d, 0xb4, 0x1f, 0x0b, 0xf6, 0x5b, 0xfb, 0xc1, 0xfb, 0x01, 0xc7, 0x1a, 0x47, 0x2f,
	0x98, 0xb8, 0xb9, 0xdd, 0x4d, 0xb8, 0xfb, 0x03, 0xf7, 0xbb, 0x6f, 0x77, 0x93, 0x1c, 0xf1, 0x0e,
	0x3f, 0x85, 0x79, 0x33, 0x6f, 0xde, 0x9b, 0x97, 0x85, 0x1f, 0x87, 0x44, 0xdf, 0x14, 0x5b, 0x16,
	0x8b, 0x6c, 0x99, 0x25, 0xb1, 0x14, 0x4b, 0x29, 0x0a, 0x8d, 0x72, 0xa1, 0x64, 0xb9, 0xcc, 0xa5,
	0xd0, 0x0d, 0x50, 0x7f, 0x98, 0xc3, 0xa2, 0x9f, 0x2f, 0x68, 0x07, 0xb1, 0xc8, 0x53, 0xae, 0xf7,
	0x42, 0x66, 0x0d, 0xa3, 0x4d, 0xaf, 0x78, 0x74, 0x0d, 0xfe, 0xdf, 0x24, 0x35, 0x35, 0xf9, 0x08,
	0x83, 0x12, 0xa5, 0x4a, 0xc4, 0x69, 0xda, 0x9d, 0x75, 0xe7, 0xc1, 0x2a, 0x60, 0xeb, 0xbb, 0x5c,
	0xa2, 0xb2, 0x10, 0xf9, 0x04, 0xc3, 0x0c, 0x35, 0xdf, 0x71, 0xcd, 0xa7, 0xaf, 0x66, 0xbd, 0x67,
	0x6d, 0xfa, 0x1b, 0xa0, 0x35, 0x1c, 0x40, 0xef, 0x88, 0xf7, 0x6e, 0xcd, 0x6b, 0x12, 0x82, 0x57,
	0xf2, 0xb4, 0x40, 0x43, 0xb3, 0xe5, 0x1b, 0x18, 0x8a, 0x1c, 0x25, 0xd7, 0x42, 0x4e, 0x7b, 0x16,
	0xa1, 0xdf, 0x60, 0x74, 0xa5, 0xb9, 0x56, 0x1b, 0xbc, 0x2d, 0x50, 0x69, 0x32, 0x81, 0x81, 0x42,
	0x59, 0x26, 0x31, 0xd6, 0x1b, 0x0c, 0x70, 0x12, 0x3b, 0xbc, 0x4e, 0x76, 0xd5, 0x0e, 0xfa, 0x05,
	0xc2, 0x9a, 0xa1, 0x72, 0x71, 0x52, 0x48, 0xde, 0x83, 0xa7, 0x2c, 0x60, 0x08, 0xd6, 0x9a, 0xcf,
	0x5c, 0x9b, 0xfe, 0x32, 0x73, 0x98, 0x62, 0xac, 0x2f, 0xae, 0xfe, 0x00, 0xfe, 0xde, 0x9d, 0x5f,
	0x1f, 0x35, 0x60, 0x55, 0x1a, 0x46, 0x62, 0xdc, 0x50, 0x6b, 0x8d, 0x77, 0xe0, 0x59, 0x17, 0x8d,
	0x86, 0xc7, 0xfe, 0x9b, 0x8a, 0xfe, 0x81, 0xe0, 0x1f, 0x97, 0xc7, 0x8b, 0x02, 0x6f, 0xa1, 0x6f,
	0x59, 0xce, 0x78, 0x43, 0xb2, 0x91, 0xa0, 0x94, 0x4f, 0x01, 0x8c, 0x61, 0x54, 0xed, 0xa8, 0x94,
	0x56, 0x0f, 0x5d, 0xf0, 0x37, 0xee, 0x27, 0x91, 0x39, 0x78, 0xee, 0x14, 0x12, 0xb2, 0x76, 0x46,
	0xd1, 0x98, 0x9d, 0x05, 0x40, 0x3b, 0xe4, 0x2b, 0xf8, 0x95, 0x61, 0x62, 0x7a, 0xed, 0xa3, 0xa3,
	0x09, 0x3b, 0xbf, 0xc4, 0x0c, 0x7f, 0x86, 0xbe, 0x55, 0x24, 0x23, 0xd6, 0x32, 0x1f, 0x85, 0xac,
	0x6d, 0x83, 0x76, 0xb6, 0xbe, 0x7b, 0x23, 0xdf, 0x1f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x88,
	0x9d, 0x0c, 0x94, 0x02, 0x00, 0x00,
}
