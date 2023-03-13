#!/usr/bin/env bash

# This script stops the blockchain created by the start_blockchain.sh script

# stop processes
pkill memcoin

# remove blockchain nodes
rm -rf /tmp/blockchain*

# remove blockchain pk
rm private.key

# kill tmux session
tmux kill-session -t blockchain
