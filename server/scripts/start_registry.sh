#!/usr/bin/env bash

# This script creates a registration server.

# Requirements:
# from server/registration: go install
# sudo apt install tmux ripgrep

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color
L=info                # default trace level
S=hbt                 # session name
W=registry            # window name
P=12000               # base port number

echo -e "${GREEN}[PARSE parameters]${NC}"
while getopts n:p:s:t:w: flag
do
    case "${flag}" in
      # t : trace level (info, debug, ...)
        t) L=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done

MASTERPANE="${S}:${W}.0"
SLAVEPANE="${S}:${W}.1"

echo -e "Split ${W} window"
tmux select-window -t ${S}:${W}
tmux splitw
tmux select-layout main-vertical

# Start the registry server
echo -e "${GREEN}[REGISTRY SERVER]${NC} starting"
tmux send-keys -t "${SLAVEPANE}" "registrycli" C-m
echo -e "${GREEN}[REGISTRY SERVER]${NC} started"

tmux select-pane -t "${MASTERPANE}"
