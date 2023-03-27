#!/usr/bin/env bash

# This script creates a new SMC (secret management committee) with a DKG.

# Requirements:
# from server/smc/smccli: go install
# sudo apt install tmux ripgrep

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color
L=info                # default trace level
S=smc                 # session name
W=smcw                # window name

echo -e "${GREEN}[PARSE parameters]${NC}"
N=4
P=11000
while getopts n:p:t: flag
do
    case "${flag}" in
      # n : number of nodes
        n) N=${OPTARG};;
      # p : starting port number
        p) P=${OPTARG};;
      # t : trace level (info, debug, ...)
        t) L=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done


echo -e "${GREEN}[TMUX setup]${NC}"
set -o errexit

# Is tmux available ?
command -v tmux >/dev/null 2>&1 || { echo >&2 "tmux is not on your PATH!"; exit 1; }

# Launch tmux session
tmux list-sessions | rg "^${S}" >/dev/null 2>&1 && { echo -e ${RED}"A session with the same name (${S}) already exists and will be destroyed${NC}"; tmux kill-session -t ${S};}


echo -e "Create tmux detached session: ${S}"
tmux new -s ${S} -n ${W} -d


echo -e "Split tmux window"
# Panes used for smc nodes
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
    p=$((P + i))
    echo -e "${GREEN}creating node #${i} on port ${p}${NC}"
    # session s, window 0, panes 1 to N
    tmux send-keys -t ${S}:${W}.${i} "LLVL=${L} LOGF=./${S}${i}.log smccli --config /tmp/${S}${i} start --listen tcp://127.0.0.1:${p}" C-m
    sleep 0.5
    i=$((i + 1));
done


echo -e "${GREEN}[CONNECT]${NC} ${N} nodes and exchange certificates"
i=2;
p=$((P + 1))
while [ ${i} -le ${N} ]
do
    # sent to master pane
    tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${S}${i} minogrpc join --address //127.0.0.1:${p} $(smccli --config /tmp/${S}1 minogrpc token)" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[INITIALIZE DKG]${NC} on each node"
i=1;
while [ ${i} -le ${N} ]
do
    # sent to master pane
    tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${S}${i} dkg listen" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[SETUP DKG]${NC} ${N} nodes"
i=1;
a=""
while [ ${i} -le ${N} ]
do
    a="${a} --authority \$(cat /tmp/${S}${i}/dkgauthority)"
    i=$((i + 1));
done
# sent to master pane
tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${S}1 dkg setup ${a} | tee smc-key.pub" C-m


# Publish the roster
echo -e "${GREEN}[PUBLISH]${NC} roster on the blockchain"
i=1;
p=$((P + i))
V="127.0.0.1:${p}"
while [ ${i} -lt ${N} ]
do
    i=$((i + 1));
    p=$((P + i));
    V="${V},127.0.0.1:${p}";
done
# Translate value to base64
#V=$(echo ${V} | base64)
K=$(grep 'Pubkey:' smc-key.pub | sed 's/🔑 Pubkey: //')

# sent to master pane
tmux send-keys -t ${S}:${W}.0 "LLVL=${L} chaincli --config /tmp/chain1 pool add\
 --key private.key\
 --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Value\
 --args value:key --args \"${K}\"\
 --args value:value --args \"${V}\"\
 --args value:command --args WRITE" C-m

# select master on pane 0
tmux select-pane -t ${S}:${W}.0

# attach to tmux session
tmux a -t ${S}
