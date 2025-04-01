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
	commonv1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/common/v1"
	serverv1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1/serverv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
)

var _ clientv1connect.GeneralsServiceHandler = &GeneralServer{}

var commands map[int][]*commonv1.Command

type GeneralServer struct {
	Name       string
	ID         string
	Log        *log.Logger
	MemberList *memberlist.Memberlist
	Commands   []string
	IsTraitor  bool
}

func (s *GeneralServer) Reset(ctx context.Context, req *connect.Request[commonv1.EmptyRequest]) (*connect.Response[commonv1.EmptyResponse], error) {
	s.Log.Info("Resetting state")

	commands = nil

	resp := commonv1.EmptyResponse{}
	return &connect.Response[commonv1.EmptyResponse]{Msg: &resp}, nil
}

func (s *GeneralServer) ReceiveCommand(ctx context.Context, req *connect.Request[v1.ReceiveCommandRequest]) (*connect.Response[commonv1.EmptyResponse], error) {
	s.Log.Info("Received command", "from", req.Msg.Command.From, "command", req.Msg.Command, "round", req.Msg.Command.Round)

	if commands == nil {
		commands = map[int][]*commonv1.Command{}
	}

	round := int(req.Msg.Command.Round)
	commands[round] = append(commands[round], req.Msg.Command)

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

	// handle the message if from the commander
	if req.Msg.Command.IsCommander && req.Msg.Command.Round == 1 {
		s.doRound1(ctx, round, commanderAddr, commanderPort)
	}

	resp := commonv1.EmptyResponse{}
	return &connect.Response[commonv1.EmptyResponse]{Msg: &resp}, nil
}

func (s *GeneralServer) doRound1(ctx context.Context, round int, commanderAddr string, commanderPort int) {
	server := serverv1connect.NewCommanderServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
	)

	r1cm := []*commonv1.Command{}
	for _, c := range commands[round] {
		r1cm = append(r1cm, c)
	}

	// the decision for round 1 is always the same as the command sent
	// by the commander
	server.DecisionMade(ctx, &connect.Request[serverv1.Decision]{
		Msg: &serverv1.Decision{
			Round:    int32(round),
			From:     s.ID,
			Decision: commands[round][0].Commands["0"],
			Commands: r1cm,
		}})
}

//func (s*GeneralServer) doRound2(ctx context.Context, round int) {
//	generals := 0
//	if req.Msg.IsCommander {
//		// send the command to all other generals
//		for _, m := range s.MemberList.Members() {
//			meta := memlist.MetaFromJSON(m.Meta)
//
//			// if the node is the commander or us then skip
//			if meta.ID == s.ID || meta.IsCommander {
//				continue
//			}
//
//			command := req.Msg.Command
//			if s.IsTraitor {
//				// if the node is a traitor then send a random command
//				command = s.Commands[generals]
//			}
//
//			s.Log.Info("Sending command", "to", meta.ID, "command", command)
//
//			client := clientv1connect.NewGeneralsServiceClient(
//				http.DefaultClient,
//				fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
//			)
//
//			// forward the command to the other node
//			client.ReceiveCommand(ctx, &connect.Request[v1.ReceiveCommandRequest]{
//				Msg: &v1.ReceiveCommandRequest{
//					Command:     command,
//					From:        s.ID,
//					IsCommander: false,
//					Round:       int32(round),
//				}})
//
//			// also send what we sent to the other node to the commander
//			// this allows the commander to build the ui
//			server := serverv1connect.NewCommanderServiceClient(
//				http.DefaultClient,
//				fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
//			)
//
//			server.CommandSent(ctx, &connect.Request[serverv1.CommandSentRequest]{
//				Msg: &serverv1.CommandSentRequest{
//					Command: command,
//					From:    s.ID,
//					To:      meta.ID,
//					Round:   int32(round),
//				}})
//
//			generals++
//		}
//
//	if len(commands[round]) == len(s.MemberList.Members())-1 {
//		decision := s.calculateDecision(round)
//		s.Log.Info("Decision made", "decision", decision)
//
//		// send the decision to the commander
//		server := serverv1connect.NewCommanderServiceClient(
//			http.DefaultClient,
//			fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
//		)
//
//		cm := []*serverv1.Command{}
//		for _, c := range commands[round] {
//			cm = append(cm, &serverv1.Command{
//				From:    c.From,
//				Command: c.Command,
//			})
//		}
//
//		server.DecisionMade(ctx, &connect.Request[serverv1.Decision]{
//			Msg: &serverv1.Decision{
//				Round:    int32(round),
//				From:     s.ID,
//				Decision: decision,
//				Commands: cm,
//			}})
//	}
//}

//func (s *GeneralServer) calculateDecision(round int) string {
//	attackCount := 0
//	retreatCount := 0
//
//	decision := "retreat"
//
//	for _, c := range commands[round] {
//		switch c.Command {
//		case "attack":
//			attackCount++
//		case "retreat":
//			retreatCount++
//		}
//	}
//
//	if attackCount > retreatCount {
//		decision = "attack"
//	}
//
//	return decision
//}
