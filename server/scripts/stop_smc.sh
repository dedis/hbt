#!/usr/bin/env bash

# This script stops the blockchain created by the start_blockchain.sh script

S=smc                 # session name

# stop processes
pkill dkgcli

# remove blockchain nodes
rm -rf /tmp/${S}*

# kill tmux session
tmux kill-session -t ${S}
