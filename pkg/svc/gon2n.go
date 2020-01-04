package svc

//go:generate sh -c "mkdir -p ../proto/generated && protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../ ../proto/*.proto"

import (
	"context"
	gon2n "github.com/pojntfx/gon2n/pkg/proto/generated/proto"
	"github.com/pojntfx/gon2n/pkg/workers"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/bloom42/libs/rz-go/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SupernodeManager manages supernodes.
type SupernodeManager struct {
	gon2n.UnimplementedSupernodeManagerServer
	SupernodesManaged map[string]*workers.Supernode
}

// Create creates a supernode.
func (s *SupernodeManager) Create(_ context.Context, args *gon2n.SupernodeManagerCreateArgs) (*gon2n.SupernodeManagerCreateReply, error) {
	id := uuid.NewV4().String()

	supernode := workers.Supernode{
		ListenPort:     int(args.GetListenPort()),
		ManagementPort: int(args.GetManagementPort()),
	}

	if err := supernode.Configure(); err != nil {
		log.Error(err.Error())

		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	if err := supernode.OpenListenPortSocket(); err != nil {
		log.Error(err.Error())

		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	if err := supernode.OpenManagementPortSocket(); err != nil {
		log.Error(err.Error())

		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	log.Info("Starting supernode")

	go supernode.Start()

	s.SupernodesManaged[id] = &supernode

	return &gon2n.SupernodeManagerCreateReply{
		Id: id,
	}, nil
}