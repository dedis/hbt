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
tmux new-session -s ${S} -n chain -d
echo -e "Create a tmux window 'smc' in the session ${S}"
tmux new-window -t ${S} -n smc

./start_chain.sh -t ${L} -s ${S} -w chain
./start_smc.sh -t ${L} -s ${S} -w smc
tmux send-keys -t ${S}:smc.0 "./publish_roster.sh" C-m

# attach to session
tmux select-pane -t ${S}:smc.0
tmux send-keys -t ${S}:smc.0 "# TMUX MINI CHEAT SHEET" C-m
tmux send-keys -t ${S}:smc.0 "# Use 'tmux lscm' to list tmux commands" C-m
tmux send-keys -t ${S}:smc.0 "# Use 'Ctrl+B N (or P)' for next (previous) window" C-m
tmux send-keys -t ${S}:smc.0 "# Use 'Ctrl+B <arrow>' to select pane" C-m
tmux send-keys -t ${S}:smc.0 "# './teardown.sh' to clean this tmux session" C-m
tmux attach -t ${S}
