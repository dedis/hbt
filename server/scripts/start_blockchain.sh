#!/usr/bin/env bash

# This script is creating a new chain and setting up the services needed to run
# an evoting system. It ends by starting the http server needed by the frontend
# to communicate with the blockchain. This operation is blocking.
# Requirements:
# sudo apt install tmux ripgrep

set -e

GREEN='\033[0;32m'
RED='\033[1;31;46m'
NC='\033[0m' # No Color

echo -e "${GREEN}[PARSE parameters]${NC}"
N=3
P=12300
while getopts n:p: flag
do
    case "${flag}" in
      # n : number of nodes
        n) N=${OPTARG};;
      # p : starting port number
        p) P=${OPTARG};;
    esac
done

echo -e "${GREEN}[TMUX setup]${NC}"
set -o errexit

# Is tmux available ?
command -v tmux >/dev/null 2>&1 || { echo >&2 "tmux is not on your PATH!"; exit 1; }

# Launch tmux session
s="blockchain"
tmux list-sessions | rg "^${s}" >/dev/null 2>&1 && { echo -e ${RED}"A session with the same name (${s})already exists and will be destroyed${NC}"; tmux kill-session -t ${s};}

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

# Start a node in each pane but the main pane
echo -e "${GREEN}[CREATE]${NC} ${N} nodes"
i=1;
while [ ${i} -le ${N} ]
do
    p=$((P+i))
    # session s, window 0, panes 1 to N
    tmux send-keys -t ${s}:0.%${i} "LLVL=info memcoin --config /tmp/node${i} start --listen tcp://127.0.0.1:${p}" C-m
    i=$((i + 1));
done

# select master on pane 0
tmux select-pane -t 0
tmux a




# Now we have N nodes running. To be continued....
exit 1

master="tmux send-keys -t $s:0.%0"
$master "./testsetup.sh" C-m

echo -e "${GREEN}[CONNECT]${NC} connect ${count}nodes"
#master="tmux send-keys -t $s:0.%0"
memcoin --config /tmp/node2 minogrpc join \
    --address //127.0.0.1:2001 $(memcoin --config /tmp/node1 minogrpc token)
memcoin --config /tmp/node3 minogrpc join \
    --address //127.0.0.1:2001 $(memcoin --config /tmp/node1 minogrpc token)
memcoin --config /tmp/node4 minogrpc join \
    --address //127.0.0.1:2001 $(memcoin --config /tmp/node1 minogrpc token)

echo -e "${GREEN}[CHAIN]${NC} create a chain"
memcoin --config /tmp/node1 ordering setup\
    --member $(memcoin --config /tmp/node1 ordering export)\
    --member $(memcoin --config /tmp/node2 ordering export)\
    --member $(memcoin --config /tmp/node3 ordering export)\
    --member $(memcoin --config /tmp/node4 ordering export)

echo -e "${GREEN}[ACCESS]${NC} setup access rights on each node"
memcoin --config /tmp/node1 access add \
    --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)
memcoin --config /tmp/node2 access add \
    --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)
memcoin --config /tmp/node3 access add \
    --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)
memcoin --config /tmp/node4 access add \
    --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)

echo -e "${GREEN}[GRANT]${NC} grant access node 1 on the chain"
memcoin --config /tmp/node1 pool add\
    --key private.key\
    --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Access\
    --args access:grant_id --args 0200000000000000000000000000000000000000000000000000000000000000\
    --args access:grant_contract --args go.dedis.ch/dela.Value\
    --args access:grant_command --args all\
    --args access:identity --args $(crypto bls signer read --path private.key --format BASE64_PUBKEY)\
    --args access:command --args GRANT
