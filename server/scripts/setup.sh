#!/usr/bin/env bash

# This script creates a new tmux session with 2 embedded windows and several panes

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color
L=info                # default trace level
S=hbt                 # session name

echo -e "${GREEN}[TMUX setup]${NC}"
set -o errexit

# Is tmux available ?
command -v tmux >/dev/null 2>&1 || { echo >&2 "tmux is not on your PATH!"; exit 1; }

# Launch tmux session
tmux list-sessions 2>/dev/null | rg "^${S}" && { echo -e ${RED}"A session with the same name (${S}) already exists and will be destroyed${NC}"; tmux kill-session -t ${S};}

echo -e "Create a tmux detached session: ${S} with a window 'chain'"
tmux new -s ${S} -n chain -d
echo -e "Create a tmux window 'smc' in the session ${S}"
tmux neww -t ${S} -n smc

./start_chain.sh -s ${S}
./start_smc.sh -s ${S}
tmux send-keys -t ${S}:smc.0 "./publish_roster.sh -s ${S} -w smc" C-m

# attach to session
tmux select-pane -t ${S}:smc.0
tmux a -t ${S}
