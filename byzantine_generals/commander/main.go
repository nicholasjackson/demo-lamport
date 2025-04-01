package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/hashicorp/memberlist"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/server/v1/serverv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/utils/env"
)

var ml *memberlist.Memberlist

func main() {
	traitor, err := env.GetBool("TRAITOR", false)
	if err != nil {
		log.Error("Failed to get traitor", "error", err)
		os.Exit(1)
	}

	commands := env.GetString("COMMANDS", "")
	arrCommands := strings.Split(commands, ",")

	grpcport, err := env.GetInt("GRPC_PORT", 8080)
	if err != nil {
		log.Error("Failed to get port", "error", err)
		os.Exit(1)
	}

	bindAddress := env.GetString("BIND_ADDR", "127.0.0.1")
	if err != nil {
		log.Error("Failed to get name", "error", err)
		os.Exit(1)
	}

	log.Info("Listening on", "memberlist", 7946, "grpc", grpcport, "name", "Commander")

	// start memberlist
	config := memberlist.DefaultLocalConfig()
	config.Name = "Commander"
	config.BindPort = 7946
	config.AdvertisePort = 7946
	config.AdvertiseAddr = bindAddress
	config.BindAddr = bindAddress
	config.Logger = log.StandardLog(log.StandardLogOptions{ForceLevel: log.DebugLevel})
	config.Delegate = &memlist.MemberListDelegate{Meta: &memlist.Meta{
		BindAddr:       bindAddress,
		MemberlistPort: 7946,
		GRPCPort:       8080,
		IsCommander:    true,
		ID:             "0",
		IsTraitor:      traitor,
	}}

	ml, err = memberlist.Create(config)
	if err != nil {
		log.Error("Failed to create memberlist", "error", err)
		os.Exit(1)
	}

	commander := &CommanderServer{
		Log:        log.NewWithOptions(os.Stderr, log.Options{Prefix: "commander"}),
		MemberList: ml,
		IsTraitor:  traitor,
		Commands:   arrCommands,
		Name:       "Commander",
		ID:         "0",
	}

	mux := http.NewServeMux()
	path, handler := serverv1connect.NewCommanderServiceHandler(commander)
	mux.Handle(path, handler)
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", bindAddress, grpcport),
		h2c.NewHandler(utils.WithCORS(mux), &http2.Server{}))
}
