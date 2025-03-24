package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/charmbracelet/log"
	"github.com/hashicorp/memberlist"
	v1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1/clientv1connect"
	serverv1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1/serverv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
)

var _ clientv1connect.GeneralsServiceHandler = &GeneralServer{}

type Command struct {
	From        string
	Command     string
	IsCommander bool
}

var commands map[int][]Command

type GeneralServer struct {
	Name       string
	ID         string
	Log        *log.Logger
	MemberList *memberlist.Memberlist
}

func (s *GeneralServer) ReceiveCommand(ctx context.Context, req *connect.Request[v1.ReceiveCommandRequest]) (*connect.Response[v1.EmptyResponse], error) {
	s.Log.Info("Received command", "from", req.Msg.From, "command", req.Msg.Command)

	if commands == nil {
		commands = map[int][]Command{
			1: []Command{},
			2: []Command{},
		}
	}

	commands[1] = append(commands[1], Command{
		From:        req.Msg.From,
		Command:     req.Msg.Command,
		IsCommander: req.Msg.IsCommander,
	})

	// if not the commander then return
	if !req.Msg.IsCommander {
		resp := v1.EmptyResponse{}
		return &connect.Response[v1.EmptyResponse]{Msg: &resp}, nil
	}

	commanderAddr := ""
	commanderPort := 0

	// grab the commander address and port
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)
		if meta.IsCommander {
			commanderAddr = meta.BindAddr
			commanderPort = meta.GRPCPort
		}
	}

	// send the command to all other generals
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)

		// if the node is the commander or us then skip
		if meta.ID == s.ID || meta.IsCommander {
			continue
		}

		s.Log.Info("Sending command", "to", meta.ID, "command", req.Msg.Command)

		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		// forward the command to the other node
		client.ReceiveCommand(ctx, &connect.Request[v1.ReceiveCommandRequest]{
			Msg: &v1.ReceiveCommandRequest{
				Command:     req.Msg.Command,
				From:        s.ID,
				IsCommander: false,
			}})

		// also send what we sent to the other node to the commander
		// this allows the commander to build the ui
		server := serverv1connect.NewCommanderServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
		)

		server.CommandSent(ctx, &connect.Request[serverv1.CommandSentRequest]{
			Msg: &serverv1.CommandSentRequest{
				Command: req.Msg.Command,
				From:    s.ID,
				To:      meta.ID,
			}})

	}

	resp := v1.EmptyResponse{}
	return &connect.Response[v1.EmptyResponse]{Msg: &resp}, nil
}
