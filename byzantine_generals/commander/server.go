package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/charmbracelet/log"
	"github.com/hashicorp/memberlist"
	clientv1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1/clientv1connect"
	v1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1/serverv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
)

var _ serverv1connect.CommanderServiceHandler = &CommanderServer{}
var nodes = []*v1.Node{
	{
		Id:             "0",
		Type:           "input",
		Data:           &v1.Data{Label: "Commander"},
		SourcePosition: "",
		TargetPosition: "",
	},
}

var commands []*v1.CommandSentRequest

var decisions []*v1.Decision

func nodeExists(id string) bool {
	for _, n := range nodes {
		if n.Id == id {
			return true
		}
	}

	return false
}

type CommanderServer struct {
	Log        *log.Logger
	MemberList *memberlist.Memberlist
}

func (s *CommanderServer) IssueCommand(ctx context.Context, req *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.CommandResponse], error) {
	s.Log.Info("Received a request to issue a command")

	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)
		if meta.IsCommander {
			continue
		}

		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		_, err := client.ReceiveCommand(ctx, &connect.Request[clientv1.ReceiveCommandRequest]{
			Msg: &clientv1.ReceiveCommandRequest{
				Command:     "attack",
				From:        "Commander",
				IsCommander: true,
			}})

		if err != nil {
			s.Log.Error("Failed to send command", "error", err)
			return nil, err
		}

		// keep a log of commands sent so we can build the edges
		commandSent := &v1.CommandSentRequest{
			Command: "attack",
			From:    "0",
			To:      meta.ID,
		}

		commands = append(commands, commandSent)
	}

	resp := v1.CommandResponse{}
	return &connect.Response[v1.CommandResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) CommandSent(ctx context.Context, req *connect.Request[v1.CommandSentRequest]) (*connect.Response[v1.CommandResponse], error) {
	s.Log.Info("Received command sent from a node")
	if commands == nil {
		commands = []*v1.CommandSentRequest{}
	}

	// add the command to the list
	commands = append(commands, req.Msg)

	resp := v1.CommandResponse{}
	return &connect.Response[v1.CommandResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) DecisionMade(ctx context.Context, req *connect.Request[v1.Decision]) (*connect.Response[v1.EmptyResponse], error) {
	s.Log.Info("Received a decision from a node", "decision", req.Msg.Decision, "from", req.Msg.From)

	decisions = append(decisions, req.Msg)

	resp := v1.EmptyResponse{}
	return &connect.Response[v1.EmptyResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) Nodes(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.NodesResponse], error) {
	s.Log.Info("Received a request for the nodes")

	// loop through the nodes and send them back
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)

		if !nodeExists(meta.ID) {
			node := &v1.Node{
				Id:             meta.ID,
				Type:           "bidirectional",
				Data:           &v1.Data{Label: m.Name},
				SourcePosition: "right",
				TargetPosition: "left",
			}

			nodes = append(nodes, node)
		}
	}

	nr := &v1.NodesResponse{
		Nodes: nodes,
	}

	resp := connect.NewResponse(nr)
	return resp, nil
}

func (s *CommanderServer) Edges(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EdgesResponse], error) {
	s.Log.Info("Received a request for the edges")

	edges := []*v1.Edge{}

	for _, e := range commands {
		edge := &v1.Edge{
			Id:     fmt.Sprintf("%s-%s", e.From, e.To),
			Source: e.From,
			Target: e.To,
			Label:  "attack",
		}

		edges = append(edges, edge)
	}

	er := &v1.EdgesResponse{
		Edges: edges,
	}

	resp := connect.NewResponse(er)
	return resp, nil
}

func (s *CommanderServer) Decisions(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.DecisionsResponse], error) {
	s.Log.Info("Received a request for the decisions")

	dr := &v1.DecisionsResponse{
		Decisions: decisions,
	}

	resp := connect.NewResponse(dr)
	return resp, nil
}
