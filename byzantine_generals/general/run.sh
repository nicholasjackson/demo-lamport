#!/bin/bash

MEMBERLIST_PORT=7951 GRPC_PORT=8085 NAME=General5 TRAITOR=true COMMANDS=retreat,retreat,attack,attack,retreat go run ./ &
general5_pid=$!

MEMBERLIST_PORT=7948 GRPC_PORT=8082 NAME=George TRAITOR=false go run ./ &
george_pid=$!

MEMBERLIST_PORT=7947 GRPC_PORT=8081 NAME=Harrald TRAITOR=false go run ./ &
harrald_pid=$!

MEMBERLIST_PORT=7949 GRPC_PORT=8083 NAME=John TRAITOR=false go run ./ &
john_pid=$!

MEMBERLIST_PORT=7950 GRPC_PORT=8084 NAME=Zoe TRAITOR=false go run ./ &
zoe_pid=$!

MEMBERLIST_PORT=7952 GRPC_PORT=8086 NAME=General6 go run ./ &
general6_pid=$!

trap ctrl_c INT

function ctrl_c() {
  kill -9 $john_pid
  kill -9 $harrald_pid
  kill -9 $george_pid
  kill -9 $zoe
  kill -9 $general5_pid
  kill -9 $general6_pid
}

# Wait for the process to finish
sleep 24h