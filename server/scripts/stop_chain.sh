#!/usr/bin/env bash

# This script stops the blockchain started by the start_chain.sh script

# stop processes
pkill chaincli

# cleanup blockchain temporary files
rm -rf /tmp/blockchain*

# remove blockchain pk
rm -f /tmp/priv.key
