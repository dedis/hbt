#!/usr/bin/env bash

# This script stops the SMC (secret management committee) created by the start_smc.sh script

S=smc                 # session name

# stop processes
pkill smccli

# remove blockchain nodes
rm -rf /tmp/${S}*

# kill tmux session
tmux kill-session -t ${S}
