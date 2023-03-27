#!/usr/bin/env bash

# This script stops the SMC (secret management committee) created by the start_smc.sh script

S=smc                 # session name

# stop processes
pkill smccli

# remove blockchain nodes
rm -rf /tmp/${S}*

# remove local public key
rm key.pub

# kill tmux session
tmux kill-session -t ${S}
