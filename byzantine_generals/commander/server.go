package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"

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

var commands []*v1.CommandSentRequest
var decisions []*v1.Decision
var votingRound = 0

type CommanderServer struct {
	Log        *log.Logger
	MemberList *memberlist.Memberlist
	IsTraitor  bool
	Commands   []string
}

func (s *CommanderServer) Reset(ctx context.Context, req *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EmptyResponse], error) {
	s.Log.Info("Resetting state")

	commands = nil
	decisions = nil
	votingRound = 0

	// reset the nodes
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)
		if meta.IsCommander {
			continue
		}

		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		client.Reset(ctx, &connect.Request[clientv1.EmptyRequest]{
			Msg: &clientv1.EmptyRequest{},
		})
	}

	resp := v1.EmptyResponse{}
	return &connect.Response[v1.EmptyResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) IssueCommand(ctx context.Context, req *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.CommandResponse], error) {
	votingRound++
	s.Log.Info("Received a request to issue a command", "round", votingRound)

	commandsSent := 0
	generals := []*memlist.Meta{}

	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)
		if meta.IsCommander {
			continue
		}

		generals = append(generals, meta)
	}

	sort.Slice(generals, func(i, j int) bool {
		return generals[i].Name < generals[j].Name
	})

	for _, meta := range generals {
		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		command := s.Commands[commandsSent]

		_, err := client.ReceiveCommand(ctx, &connect.Request[clientv1.ReceiveCommandRequest]{
			Msg: &clientv1.ReceiveCommandRequest{
				Command:     command,
				From:        "Commander",
				IsCommander: true,
				Round:       int32(votingRound),
			}})

		if err != nil {
			s.Log.Error("Failed to send command", "error", err)
			return nil, err
		}

		// keep a log of commands sent so we can build the edges
		commandSent := &v1.CommandSentRequest{
			Command: command,
			From:    "0",
			To:      meta.ID,
		}

		s.Log.Info("Sending command", "to", meta.ID, "command", command)

		commandsSent++
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
	nodes := []*v1.Node{}

	// loop through the nodes and send them back
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)

		nodeType := "bidirectional"
		nodeSourcePosition := "right"
		nodeTargetPosition := "left"

		if meta.IsCommander {
			nodeType = "input"
			nodeSourcePosition = ""
			nodeTargetPosition = ""
		}

		node := &v1.Node{
			Id:             meta.ID,
			Type:           nodeType,
			Data:           &v1.Data{Label: m.Name},
			SourcePosition: nodeSourcePosition,
			TargetPosition: nodeTargetPosition,
			IsTraitor:      meta.IsTraitor,
		}

		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Data.Label < nodes[j].Data.Label
	})

	nr := &v1.NodesResponse{
		Nodes: nodes,
	}

	resp := connect.NewResponse(nr)
	return resp, nil
}

func (s *CommanderServer) Edges(context.Context, *connect.Request[v1.EmptyRequest]) (*connect.Response[v1.EdgesResponse], error) {
	edges := []*v1.Edge{}

	for _, e := range commands {
		if e.Round != int32(votingRound) {
			continue
		}

		edge := &v1.Edge{
			Id:     fmt.Sprintf("%s-%s", e.From, e.To),
			Source: e.From,
			Target: e.To,
			Label:  e.Command,
		}

		edges = append(edges, edge)
	}

	er := &v1.EdgesResponse{
		Edges: edges,
	}

	resp := connect.NewResponse(er)
	return resp, nil
}

func (s *CommanderServer) Decisions(ctx context.Context, r *connect.Request[v1.DecisionsRequest]) (*connect.Response[v1.DecisionsResponse], error) {
	ds := []*v1.Decision{}

	for _, d := range decisions {
		if d.Round == int32(votingRound) || r.Msg.AllData {
			ds = append(ds, d)
		}
	}

	dr := &v1.DecisionsResponse{
		Decisions: ds,
	}

	resp := connect.NewResponse(dr)
	return resp, nil
}
