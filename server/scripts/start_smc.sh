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
S=hbt                 # session name
W=smc                 # window name
N=4                   # number of nodes
P=11000               # base port number

echo -e "${GREEN}[PARSE parameters]${NC}"
while getopts n:p:s:t:w: flag
do
    case "${flag}" in
      # n : number of nodes
        n) N=${OPTARG};;
      # p : starting port number
        p) P=${OPTARG};;
      # s : session name
        s) S=${OPTARG};;
      # t : trace level (info, debug, ...)
        t) L=${OPTARG};;
      # w : window name
        w) W=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done


echo -e "Split ${W} window"
tmux select-window -t ${S}:${W}
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
    tmux send-keys -t ${S}:${W}.${i} "LLVL=${L} LOGF=./${W}${i}.log smccli --config /tmp/${W}${i} start --listen tcp://127.0.0.1:${p}" C-m
    sleep 0.5
    i=$((i + 1));
done


echo -e "${GREEN}[CONNECT]${NC} ${N} nodes and exchange certificates"
i=2;
p=$((P + 1))
while [ ${i} -le ${N} ]
do
    tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${W}${i} minogrpc join --address //127.0.0.1:${p} $(smccli --config /tmp/${W}1 minogrpc token)" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[INITIALIZE DKG]${NC} on each node"
i=1;
while [ ${i} -le ${N} ]
do
    tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${W}${i} dkg listen" C-m
    i=$((i + 1));
done


echo -e "${GREEN}[SETUP DKG]${NC} ${N} nodes"
i=1;
a=""
while [ ${i} -le ${N} ]
do
    a="${a} --authority \$(cat /tmp/${W}${i}/dkgauthority)"
    i=$((i + 1));
done
tmux send-keys -t ${S}:${W}.0 "smccli --config /tmp/${W}1 dkg setup ${a} | tee smckey.pub" C-m


# Publish the roster
echo -e "${GREEN}[SAVE]${NC} roster to file"
i=1;
p=$((P + i))
V="127.0.0.1:${p}"
while [ ${i} -lt ${N} ]
do
    i=$((i + 1));
    p=$((P + i));
    V="${V},127.0.0.1:${p}";
done
tmux send-keys -t ${S}:${W}.0 "echo ${V} > roster.txt" C-m

tmux select-pane -t ${S}:${W}.0
