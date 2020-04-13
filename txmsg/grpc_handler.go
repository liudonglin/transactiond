package txmsg

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

var groups = make(map[string][]string)

type gRpcServer struct{}

func (s *gRpcServer) CreateGroup(ctx context.Context, in *CreateGroupMessage) (*ResponseMessage, error) {
	log.Printf("create a new group : [%s]", in.GroupId)
	groups[in.GroupId] = make([]string, 0)
	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) JoinGroup(ctx context.Context, in *JoinGroupMessage) (*ResponseMessage, error) {
	log.Printf("a unit [%s] join the group [%s]", in.UnitId, in.GroupId)
	if v, ok := groups[in.GroupId]; ok {
		v = append(v, in.UnitId)
	}
	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) NotifyGroup(ctx context.Context, in *NotifyGroupMessage) (*ResponseMessage, error) {
	log.Printf("group received notify [%s], stats is [%s] ", in.GroupId, in.State.String())
	if v, ok := groups[in.GroupId]; ok {
		print(len(v))
	}
	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) NotifyUnit(in *NotifyUnitMessage, out ManageService_NotifyUnitServer) error {
	log.Printf("group notify [%s]", in.GroupId)
	if v, ok := groups[in.GroupId]; ok {
		print(len(v))
	}
	return nil
}

func StartRpcHandler() {
	s := grpc.NewServer()
	RegisterManageServiceServer(s, &gRpcServer{})

	gRpcPort := ":6102"
	lis, err := net.Listen("tcp", gRpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
