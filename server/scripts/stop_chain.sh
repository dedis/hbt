#!/usr/bin/env bash

# This script stops the blockchain created by the start_chain.sh script

S=chain               # session name

# stop processes
pkill chaincli

# remove blockchain nodes
rm -rf /tmp/${S}*

# remove blockchain pk
rm private.key

# kill tmux session
tmux kill-session -t ${S}
