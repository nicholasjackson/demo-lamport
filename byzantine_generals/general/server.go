package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"

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
	s.Log.Debug("Received command", "name", s.Name, "from", req.Msg.Command.From, "round", req.Msg.Command.Round)
	s.Log.Debug("command", "value", req.Msg.Command)

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
	switch round {
	case 1:
		// if we are in round 1 then make a decision based on the commanders message
		if req.Msg.Command.IsCommander {
			s.decisionRound1(ctx, commanderAddr, commanderPort)
		}
	case 2:
		// If we are in round 2 we need to send the command we recieved from the commander
		// to all other generals.
		if req.Msg.Command.IsCommander {
			s.sendRound2Messages(ctx)
		}

		// Then we need to make a decision based on the messages we have received
		// from the other generals

		// This decision based on the aggregated commands from the other generals
		// not just the commander will allow us to make a decision based on the majority
		// If there is a single traitor then that traitors command will be ignored
		if len(commands[2]) == len(s.MemberList.Members())-1 {
			s.decisionRound2(ctx, commanderAddr, commanderPort)
		}
	case 3:
		// If we are in round 3 then we need to send all the commands we have received
		// from all the generals to the other generals
		// This allows us to handle two traitors
		if req.Msg.Command.IsCommander {
			s.sendRound3Messages(ctx)
		}

		// To handle two traitors we need to look at the messages that other generals
		// recieved from the other generals not just the commander
		// this allows us to filter out the traitors
		if len(commands[3]) == len(s.MemberList.Members())-1 {
			s.decisionRound3(ctx, commanderAddr, commanderPort)
		}
	}

	resp := commonv1.EmptyResponse{}
	return &connect.Response[commonv1.EmptyResponse]{Msg: &resp}, nil
}

// send the decision to the commander
func (s *GeneralServer) decisionRound1(ctx context.Context, commanderAddr string, commanderPort int) {
	server := serverv1connect.NewCommanderServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
	)

	r1cm := []*commonv1.Command{}
	r1cm = append(r1cm, commands[1]...)

	// the decision for round 1 is always the same as the command sent
	// by the commander
	server.DecisionMade(ctx, &connect.Request[serverv1.Decision]{
		Msg: &serverv1.Decision{
			Round:    1,
			From:     s.ID,
			Decision: commands[1][0].Commands["0"],
			Commands: r1cm,
		}})

	s.Log.Info("Decision made", "round", 1, "general", s.Name, "id", s.ID, "commands", len(commands[1]), "decision", commands[1][0].Commands["0"])
}

// send the commanders message to all other generals
func (s *GeneralServer) sendRound2Messages(ctx context.Context) {
	s.Log.Info("Sending round 2 messages", "general", s.Name)

	// get the list of generals in alphabetical order
	generals := []*memlist.Meta{}
	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)
		if meta.ID == s.ID || meta.IsCommander {
			continue
		}

		generals = append(generals, meta)
	}

	sort.Slice(generals, func(i, j int) bool {
		return generals[i].Name < generals[j].Name
	})

	for i, meta := range generals {
		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		co := commands[1][0].Commands["0"]
		if s.IsTraitor {
			// if we are a traitor then we need to send a different command
			// to the other generals
			co = s.Commands[i]
			s.Log.Info("Sending traitor command", "name", s.Name, "to", meta.ID, "round", 2, "command", co)
		}

		// forward the command to the other node
		// but change the details to us
		command := &commonv1.Command{
			Commands:    map[string]string{s.ID: co},
			From:        s.ID,
			IsCommander: false,
			Round:       2,
		}

		client.ReceiveCommand(ctx, &connect.Request[v1.ReceiveCommandRequest]{
			Msg: &v1.ReceiveCommandRequest{
				Command: command,
			}})

		s.Log.Debug("Sent command", "name", s.Name, "to", meta.ID, "round", 2, "commands", command)
	}
}

func (s *GeneralServer) decisionRound2(ctx context.Context, commanderAddr string, commanderPort int) {
	// sum the commands from the generals
	attackCount := 0
	retreatCount := 0
	for _, c := range commands[2] {
		for _, v := range c.Commands {
			if v == "attack" {
				attackCount++
			} else {
				retreatCount++
			}
		}
	}

	decision := "retreat"
	if attackCount > retreatCount {
		decision = "attack"
	}

	r2cm := []*commonv1.Command{}
	r2cm = append(r2cm, commands[2]...)

	server := serverv1connect.NewCommanderServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
	)

	// the decision for round 1 is always the same as the command sent
	// by the commander
	server.DecisionMade(ctx, &connect.Request[serverv1.Decision]{
		Msg: &serverv1.Decision{
			Round:    2,
			From:     s.ID,
			Decision: decision,
			Commands: r2cm,
		}})

	s.Log.Info("Decision made", "round", 2, "general", s.Name, "commands", len(commands[2]), "decision", decision)
}

func (s *GeneralServer) sendRound3Messages(ctx context.Context) {
	// for each command we have received from the other generals
	// send it to the other node
	combined := map[string]string{}
	for _, c := range commands[2] {
		for k, v := range c.Commands {
			combined[k] = v
		}
	}

	command := &commonv1.Command{
		Commands:    combined,
		From:        s.ID,
		IsCommander: false,
		Round:       3,
	}

	for _, m := range s.MemberList.Members() {
		meta := memlist.MetaFromJSON(m.Meta)

		// if the node is the commander or us then skip
		if meta.ID == s.ID || meta.IsCommander {
			continue
		}

		client := clientv1connect.NewGeneralsServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%d", meta.BindAddr, meta.GRPCPort),
		)

		client.ReceiveCommand(ctx, &connect.Request[v1.ReceiveCommandRequest]{
			Msg: &v1.ReceiveCommandRequest{
				Command: command,
			}})

		s.Log.Debug("Sent command", "name", s.Name, "to", meta.ID, "round", 3, "command", command)
	}
}

func (s *GeneralServer) decisionRound3(ctx context.Context, commanderAddr string, commanderPort int) {
	s.Log.Info("Making decision", "general", s.Name, "round", 3, "commands", len(commands[3]))

	for _, c := range commands[3] {
		s.Log.Debug("Commands", "from", c.From, "commands", c.Commands)
	}

	// we need to count the number of attack and retreat commands from each general
	// then we can work out what the consensus decision is from each one
	generalResults := map[string][]string{}
	for _, c := range commands[3] {
		for k, v := range c.Commands {
			generalResults[k] = append(generalResults[k], v)
		}
	}

	// now sum the results for each general
	generalDecisions := map[string]string{}
	for k, v := range generalResults {
		attackCount := 0
		retreatCount := 0
		for _, d := range v {
			if d == "attack" {
				attackCount++
			} else {
				retreatCount++
			}
		}

		if attackCount > retreatCount {
			generalDecisions[k] = "attack"
		} else {
			generalDecisions[k] = "retreat"
		}
	}

	// now sum the aggregated results
	attackCount := 0
	retreatCount := 0
	for _, v := range generalDecisions {
		if v == "attack" {
			attackCount++
		} else {
			retreatCount++
		}
	}

	decision := "retreat"
	if attackCount > retreatCount {
		decision = "attack"
	}

	server := serverv1connect.NewCommanderServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:%d", commanderAddr, commanderPort),
	)

	r3cm := []*commonv1.Command{}
	r3cm = append(r3cm, commands[3]...)

	// the decision for round 1 is always the same as the command sent
	// by the commander
	server.DecisionMade(ctx, &connect.Request[serverv1.Decision]{
		Msg: &serverv1.Decision{
			Round:    3,
			From:     s.ID,
			Decision: decision,
			Commands: r3cm,
		}})

	s.Log.Info("Decision made", "general", s.Name, "round", 3, "commands", len(commands[3]), "decision", decision)
}
