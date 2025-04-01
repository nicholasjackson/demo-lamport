package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	traitor, err := env.GetBool("TRAITOR", false)
	if err != nil {
		log.Error("Failed to get port", "error", err)
		os.Exit(1)
	}

	commands := env.GetString("COMMANDS", "")
	arrCommands := strings.Split(commands, ",")

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

	debug, _ := env.GetBool("DEBUG", false)
	level := log.InfoLevel
	if debug {
		level = log.DebugLevel
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "general",
		Level:  level,
	})

	logger.Info("Listening", "memberlist", mlport, "grpc", grpcport, "name", name)
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
		IsTraitor:      traitor,
		Name:           name,
	}}

	ml, err = memberlist.Create(config)
	if err != nil {
		logger.Error("Failed to create memberlist", "error", err)
		os.Exit(1)
	}

	nodes, err := ml.Join([]string{commanderAddress})
	if err != nil {
		logger.Error("Failed to join memberlist", "error", err)
		os.Exit(1)
	}

	logger.Info("Joined memberlist", "nodes", nodes)

	general := &GeneralServer{
		Log:        logger,
		MemberList: ml,
		Name:       name,
		ID:         guid.String(),
		IsTraitor:  traitor,
		Commands:   arrCommands,
	}

	go func() {
		mux := http.NewServeMux()
		path, handler := clientv1connect.NewGeneralsServiceHandler(general)
		mux.Handle(path, handler)
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", bindAddress, grpcport),
			h2c.NewHandler(utils.WithCORS(mux), &http2.Server{}),
		)
	}()

	// block until ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Info("Exiting")

	err = ml.Leave(0)
	if err != nil {
		log.Error("Failed to leave memberlist", "error", err)
		os.Exit(1)
	}
}
