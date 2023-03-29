#!/usr/bin/env bash

# This script stops the blockchain created by the start_chain.sh script

W=chain               # window name

# stop processes
pkill chaincli

# remove blockchain nodes
rm -rf /tmp/${W}*

# remove blockchain pk
rm -f private.key
