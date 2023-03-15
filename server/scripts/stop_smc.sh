#!/usr/bin/env bash

# This script stops the blockchain created by the start_blockchain.sh script

# stop processes
pkill dkgcli

# remove blockchain nodes
rm -rf /tmp/smc*

# kill tmux session
tmux kill-session -t smc
