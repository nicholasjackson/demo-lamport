#!/bin/bash
MEMBERLIST_PORT=7948 GRPC_PORT=8082 NAME=George TRAITOR=false go run ./ &
george_pid=$!

MEMBERLIST_PORT=7947 GRPC_PORT=8081 NAME=Harrald TRAITOR=false go run ./ &
harrald_pid=$!

MEMBERLIST_PORT=7949 GRPC_PORT=8083 NAME=John TRAITOR=false go run ./ &
john_pid=$!

MEMBERLIST_PORT=7951 GRPC_PORT=8085 NAME=Leif TRAITOR=true COMMANDS=retreat,attack,attack,attack,retreat go run ./ &
leif_pid=$!

MEMBERLIST_PORT=7952 GRPC_PORT=8086 NAME=Rolo go run ./ &
rolo_pid=$!

MEMBERLIST_PORT=7950 GRPC_PORT=8084 NAME=Zoe TRAITOR=false go run ./ &
zoe_pid=$!

trap ctrl_c INT

function ctrl_c() {
  kill -9 $john_pid
  kill -9 $harrald_pid
  kill -9 $george_pid
  kill -9 $zoe
  kill -9 $leif_pid
  kill -9 $rolo_pid
}

# Wait for the process to finish
sleep 24h