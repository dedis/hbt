#!/usr/bin/env bash

# This script stops the blockchain created by the start_chain.sh script

W=chain               # window name

# stop processes
pkill chaincli

# cleanup blockchain temporary files
rm -rf /tmp/${W}*

# remove blockchain pk
rm -f private.key
