#!/usr/bin/env bash

# This script tears down the blockchain, the smc, and the tmux session

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color

S=hbt                 # session name

echo -e "${GREEN}[TMUX tear down]${NC}"

./stop_registry.sh
./stop_smc.sh
./stop_chain.sh

tmux kill-session -t ${S}

