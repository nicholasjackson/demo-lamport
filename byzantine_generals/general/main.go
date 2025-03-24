package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/charmbracelet/log"
	"github.com/hashicorp/memberlist"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/client/v1/clientv1connect"
	memlist "github.com/nicholasjackson/demo-lamport/byzantine_generals/memberlist"
	"github.com/nicholasjackson/demo-lamport/byzantine_generals/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/utils/env"
)

var ml *memberlist.Memberlist

func main() {
	mlport, err := env.GetInt("MEMBERLIST_PORT", 7947)
	if err != nil {
		log.Error("Failed to get port", "error", err)
		os.Exit(1)
	}

	grpcport, err := env.GetInt("GRPC_PORT", 8081)
	if err != nil {
		log.Error("Failed to get port", "error", err)
		os.Exit(1)
	}

	name := env.GetString("NAME", "commander")
	if err != nil {
		log.Error("Failed to get name", "error", err)
		os.Exit(1)
	}

	bindAddress := env.GetString("BIND_ADDR", "127.0.0.1")
	if err != nil {
		log.Error("Failed to get name", "error", err)
		os.Exit(1)
	}

	commanderAddress := env.GetString("COMMANDER_ADDR", "127.0.0.1:7946")
	if err != nil {
		log.Error("Failed to get name", "error", err)
		os.Exit(1)
	}

	log.Info("Listening", "memberlist", mlport, "grpc", grpcport, "name", name)
	guid, _ := guid.NewV4()

	config := memberlist.DefaultLocalConfig()
	config.Name = name
	config.BindPort = mlport
	config.AdvertisePort = mlport
	config.AdvertiseAddr = bindAddress
	config.BindAddr = bindAddress
	config.Logger = log.StandardLog(log.StandardLogOptions{ForceLevel: log.DebugLevel})
	config.Delegate = &memlist.MemberListDelegate{Meta: &memlist.Meta{
		BindAddr:       bindAddress,
		MemberlistPort: mlport,
		GRPCPort:       grpcport,
		ID:             guid.String(),
		IsCommander:    false,
		Name:           name,
	}}

	ml, err = memberlist.Create(config)
	if err != nil {
		log.Error("Failed to create memberlist", "error", err)
		os.Exit(1)
	}

	nodes, err := ml.Join([]string{commanderAddress})
	if err != nil {
		log.Error("Failed to join memberlist", "error", err)
		os.Exit(1)
	}

	log.Info("Joined memberlist", "nodes", nodes)

	general := &GeneralServer{
		Log:        log.NewWithOptions(os.Stderr, log.Options{Prefix: "general"}),
		MemberList: ml,
		Name:       name,
		ID:         guid.String(),
	}

	mux := http.NewServeMux()
	path, handler := clientv1connect.NewGeneralsServiceHandler(general)
	mux.Handle(path, handler)
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", bindAddress, grpcport),
		h2c.NewHandler(utils.WithCORS(mux), &http2.Server{}),
	)
}
