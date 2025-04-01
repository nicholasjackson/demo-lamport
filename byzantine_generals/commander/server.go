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
	commonv1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/common/v1"
	v1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1/serverv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
)

var _ serverv1connect.CommanderServiceHandler = &CommanderServer{}

var decisions []*v1.Decision
var votingRound = 0

type CommanderServer struct {
	Name       string
	ID         string
	Log        *log.Logger
	MemberList *memberlist.Memberlist
	IsTraitor  bool
	Commands   []string
}

func (s *CommanderServer) Reset(ctx context.Context, req *connect.Request[commonv1.EmptyRequest]) (*connect.Response[commonv1.EmptyResponse], error) {
	s.Log.Info("Resetting state")

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

		client.Reset(ctx, &connect.Request[commonv1.EmptyRequest]{
			Msg: &commonv1.EmptyRequest{},
		})
	}

	resp := commonv1.EmptyResponse{}
	return &connect.Response[commonv1.EmptyResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) IssueCommand(ctx context.Context, req *connect.Request[commonv1.EmptyRequest]) (*connect.Response[v1.CommandResponse], error) {
	votingRound++
	s.Log.Info("Received a request to issue a command", "round", votingRound)

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

	for i, meta := range generals {
		s.Log.Info("Sending command", "to", meta.ID, "round", votingRound, "command", s.Commands[i])

		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		command := &commonv1.Command{
			Commands:    map[string]string{s.ID: s.Commands[i]},
			From:        s.ID,
			IsCommander: true,
			Round:       int32(votingRound),
		}

		_, err := client.ReceiveCommand(
			ctx,
			&connect.Request[clientv1.ReceiveCommandRequest]{
				Msg: &clientv1.ReceiveCommandRequest{
					Command: command,
				},
			},
		)

		if err != nil {
			s.Log.Error("Failed to send command", "error", err)
			return nil, err
		}
	}

	resp := v1.CommandResponse{}
	return &connect.Response[v1.CommandResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) DecisionMade(ctx context.Context, req *connect.Request[v1.Decision]) (*connect.Response[commonv1.EmptyResponse], error) {
	s.Log.Info("Received a decision from a node", "decision", req.Msg.Decision, "from", req.Msg.From)

	decisions = append(decisions, req.Msg)

	resp := commonv1.EmptyResponse{}
	return &connect.Response[commonv1.EmptyResponse]{Msg: &resp}, nil
}

func (s *CommanderServer) Nodes(context.Context, *connect.Request[commonv1.EmptyRequest]) (*connect.Response[v1.NodesResponse], error) {
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

func (s *CommanderServer) Edges(context.Context, *connect.Request[commonv1.EmptyRequest]) (*connect.Response[v1.EdgesResponse], error) {
	edges := []*v1.Edge{}
	generals := []*memlist.Meta{}

	// if we have not voted return nothing
	if votingRound == 0 {
		er := &v1.EdgesResponse{
			Edges: edges,
		}
		resp := connect.NewResponse(er)
		return resp, nil
	}

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

	// add the edges from the commander to the generals
	for i, g := range generals {
		edge := &v1.Edge{
			Id:     fmt.Sprintf("0-%s", g.ID),
			Source: "0",
			Target: g.ID,
			Label:  s.Commands[i],
		}

		edges = append(edges, edge)
	}

	// add the edges from the generals to each other
	for _, d := range decisions {
		if d.Round != 2 {
			continue
		}

		for _, c := range d.Commands {
			if c.From == "0" {
				continue
			}

			// find the general that sent data to this general
			// and add an edge
			edge := &v1.Edge{
				Id:     fmt.Sprintf("%s-%s", c.From, d.From),
				Source: c.From,
				Target: d.From,
				Label:  c.Commands[c.From],
			}

			edges = append(edges, edge)
		}

		//ds = append(ds, d)
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
