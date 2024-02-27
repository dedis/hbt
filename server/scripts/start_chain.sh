#!/usr/bin/env bash

# This script creates a new blockchain and grants access to the first node.

# Requirements:
# from server/blockchain/chaincli: go install
# from dela/cli/crypto: go install
# sudo apt install tmux ripgrep

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color
L=info                # default trace level
S=hbt                 # session name
W=blockchain
N=4                   # number of nodes
P=10000               # base port number

echo -e "${GREEN}[PARSE parameters]${NC}"
while getopts n:p:s:t:w: flag
do
    case "${flag}" in
      # n : number of nodes
        n) N=${OPTARG};;
      # t : trace level (info, debug, ...)
        t) L=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done

# master pane
MASTERPANE="${S}:${W}.0"


echo -e "Split ${W} window"
tmux select-window -t ${S}:${W}
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
    echo -e "${GREEN}creating node #${i} on port ${p}${NC}"
    # session s, window 0, panes 1 to N
    tmux send-keys -t ${S}:${W}.${i} "LLVL=${L} LOGF=./${W}${i}.log chaincli --config /tmp/${W}${i} start --listen tcp://127.0.0.1:${p}" C-m
    sleep 0.5
    i=$((i + 1));
done


echo -e "${GREEN}[CONNECT]${NC} ${N} nodes and exchange certificates"
i=2;
p=$((P + 1))
while [ ${i} -le ${N} ]
do
    echo -e "joining node ${i} on master pane ${MASTERPANE}"
    tmux send-keys -t "${MASTERPANE}" "chaincli --config /tmp/${W}${i} minogrpc join --address //127.0.0.1:${p} $(chaincli --config /tmp/${W}1 minogrpc token)" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[CHAIN]${NC} ${N} nodes"
i=1;
m=""
while [ ${i} -le ${N} ]
do
    m="${m} --member \$(chaincli --config /tmp/${W}${i} ordering export)"
    i=$((i + 1));
done
tmux send-keys -t "${MASTERPANE}" "chaincli --config /tmp/${W}1 ordering setup ${m}" C-m


echo -e "${GREEN}[ACCESS]${NC} setup access rights on each node"
i=1;
while [ ${i} -le ${N} ]
do
    tmux send-keys -t "${MASTERPANE}" "chaincli --config /tmp/${W}${i} access add \
                                  --identity $(crypto bls signer read --path private.key --format BASE64_PUBKEY)" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[GRANT]${NC} grant access for node 1 on the chain"
# sent to master pane
tmux send-keys -t "${MASTERPANE}" "chaincli --config /tmp/${W}1 pool add\
    --key private.key\
    --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Access\
    --args access:grant_id --args 0200000000000000000000000000000000000000000000000000000000000000\
    --args access:grant_contract --args go.dedis.ch/dela.Value\
    --args access:grant_command --args all\
    --args access:identity --args $(crypto bls signer read --path private.key --format BASE64_PUBKEY)\
    --args access:command --args GRANT" C-m

tmux select-pane -t "${MASTERPANE}"
