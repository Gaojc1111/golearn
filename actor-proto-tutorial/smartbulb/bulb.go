package smartbulb

import (
	"github.com/asynkron/protoactor-go/cluster"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Bulb struct {
	on bool
}

func (b *Bulb) Init(ctx cluster.GrainContext)           {}
func (b *Bulb) Terminate(ctx cluster.GrainContext)      {}
func (b *Bulb) ReceiveDefault(ctx cluster.GrainContext) {}

func (b *Bulb) TurnOn(_ *emptypb.Empty, ctx cluster.GrainContext) (*emptypb.Empty, error) {
	b.on = true
	return &emptypb.Empty{}, nil
}

func (b *Bulb) TurnOff(_ *emptypb.Empty, ctx cluster.GrainContext) (*emptypb.Empty, error) {
	b.on = false
	return &emptypb.Empty{}, nil
}

func (b *Bulb) GetState(_ *emptypb.Empty, ctx cluster.GrainContext) (*GetStateResponse, error) {
	return &GetStateResponse{IsOn: b.on}, nil
}

func init() {
	SmartBulbFactory(func() SmartBulb { return &Bulb{} })
}
