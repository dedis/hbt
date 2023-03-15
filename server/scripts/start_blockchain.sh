#!/usr/bin/env bash

# This script is creating a new chain and setting up the services needed to run
# an evoting system. It ends by starting the http server needed by the frontend
# to communicate with the blockchain. This operation is blocking.

# Requirements:
# from dela/cli/node/memcoin: go install
# from dela/cli/crypto: go install
# sudo apt install tmux ripgrep

set -e

GREEN='\033[0;32m'
RED='\033[1;31;46m'
NC='\033[0m' # No Color

echo -e "${GREEN}[PARSE parameters]${NC}"
N=3
P=10000
while getopts n:p: flag
do
    case "${flag}" in
      # n : number of nodes
        n) N=${OPTARG};;
      # p : starting port number
        p) P=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done

echo -e "${GREEN}[TMUX setup]${NC}"
set -o errexit

# Is tmux available ?
command -v tmux >/dev/null 2>&1 || { echo >&2 "tmux is not on your PATH!"; exit 1; }

# Launch tmux session
s="blockchain"
tmux list-sessions | rg "^${s}" >/dev/null 2>&1 && { echo -e ${RED}"A session with the same name (${s}) already exists and will be destroyed${NC}"; tmux kill-session -t ${s};}

echo -e "Create tmux detached session: ${s}"
tmux new -s $s -n nodes -d

echo -e "Split tmux window"
# Panes used for blockchain nodes
i=1;
while [ ${i} -le ${N} ]
do
    tmux splitw
    tmux select-layout main-vertical
    i=$((i + 1));
done

echo -e "${GREEN}[PK]${NC} create a private key"
crypto bls signer new --save private.key

# Start a node in each pane but the main pane
echo -e "${GREEN}[CREATE]${NC} ${N} nodes"
i=1;
while [ ${i} -le ${N} ]
do
    p=$((P + i))
    # session s, window 0, panes 1 to N
    echo -e "${GREEN}creating node #${N} on port ${p}${NC}"
    tmux send-keys -t ${s}:0.%${i} "LLVL=info memcoin --config /tmp/blockchain${i} start --listen tcp://127.0.0.1:${p}" C-m
    sleep 0.5
    i=$((i + 1));
done

echo -e "${GREEN}[CONNECT]${NC} ${N} nodes"
i=2;
p=$((P + 1))
while [ ${i} -le ${N} ]
do
    # sent to master pane
    tmux send-keys -t ${s}:0.%0 "LLVL=info memcoin --config /tmp/blockchain${i} minogrpc join --address //127.0.0.1:${p} $(memcoin --config /tmp/blockchain1 minogrpc token)" C-m
    i=$((i + 1));
done

echo -e "${GREEN}[CHAIN]${NC} ${N} nodes"
# sent to master pane
# TODO: convert to a N loop
tmux send-keys -t ${s}:0.%0 "memcoin --config /tmp/blockchain1 ordering setup\
    --member $(memcoin --config /tmp/blockchain1 ordering export)\
    --member $(memcoin --config /tmp/blockchain2 ordering export)\
    --member $(memcoin --config /tmp/blockchain3 ordering export)" C-m

echo -e "${GREEN}[ACCESS]${NC} setup access rights on each node"
i=1;
while [ ${i} -le ${N} ]
do
    # sent to master pane
    tmux send-keys -t ${s}:0.%0 "LLVL=info memcoin --config /tmp/blockchain${i} access add \
                                  --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)" C-m
    i=$((i + 1));
done

echo -e "${GREEN}[GRANT]${NC} grant access for node 1 on the chain"
# sent to master pane
tmux send-keys -t ${s}:0.%0 "memcoin --config /tmp/blockchain1 pool add\
    --key private.key\
    --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Access\
    --args access:grant_id --args 0200000000000000000000000000000000000000000000000000000000000000\
    --args access:grant_contract --args go.dedis.ch/dela.Value\
    --args access:grant_command --args all\
    --args access:identity --args $(crypto bls signer read --path private.key --format BASE64_PUBKEY)\
    --args access:command --args GRANT" C-m

# select master on pane 0
tmux select-pane -t 0
tmux a
