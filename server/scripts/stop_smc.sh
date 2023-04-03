#!/usr/bin/env bash

# This script stops the SMC (secret management committee) created by the start_smc.sh script

W=smc                 # window name

# stop processes
pkill smccli

# cleanup temporary SMC files
rm -rf /tmp/${W}*

# remove local public key
rm -f smckey.pub

