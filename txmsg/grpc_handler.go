package txmsg

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"transactiond/common"
)

var groups = make(map[string][]*common.Transaction)
var clients = make(map[string]chan *NotifyUnitMessage)

type gRpcServer struct{}

func (s *gRpcServer) CreateGroup(ctx context.Context, in *CreateGroupMessage) (*ResponseMessage, error) {
	log.Printf("create a new group : [%s]", in.GroupId)
	groups[in.GroupId] = make([]*common.Transaction, 0)
	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) JoinGroup(ctx context.Context, in *JoinGroupMessage) (*ResponseMessage, error) {
	log.Printf("a unit [%s] join the group [%s]", in.UnitId, in.GroupId)

	p, _ := peer.FromContext(ctx)
	addr := p.Addr.String()

	if infos, ok := groups[in.GroupId]; ok {
		info := &common.Transaction{
			GlobalId: in.GroupId,
			UnitId:   in.UnitId,
			Source:   addr,
		}
		infos = append(infos, info)
		groups[in.GroupId] = infos
	}
	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) NotifyGroup(ctx context.Context, in *NotifyGroupMessage) (*ResponseMessage, error) {
	log.Printf("group received notify [%s], stats is [%s] ", in.GroupId, in.State.String())
	if infos, ok := groups[in.GroupId]; ok {

		msg := &NotifyUnitMessage{
			GroupId: in.GroupId,
			State:   in.State,
		}

		for _, info := range infos {
			client := clients[info.Source]
			client <- msg
		}
	}

	return &ResponseMessage{Code: "200"}, nil
}

func (s *gRpcServer) NotifyUnit(in *NotifyUnitMessage, out ManageService_NotifyUnitServer) error {
	p, _ := peer.FromContext(out.Context())
	addr := p.Addr.String()

	log.Printf("a new client connection address is [%s]", addr)

	notifyChan := make(chan *NotifyUnitMessage)
	clients[addr] = notifyChan

	for {
		select {
		case msg := <-notifyChan:
			out.Send(msg)
		}
	}
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
