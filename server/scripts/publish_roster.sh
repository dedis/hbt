#!/usr/bin/env bash

# This script publishes the SMC roster on the blockchain.

# Requirements:
# from server/blockchain/chaincli: go install

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color
S=hbt                 # session name
W=smc                 # window name

echo -e "${GREEN}[PARSE parameters]${NC}"
while getopts s:w: flag
do
    case "${flag}" in
      # s : session name
        s) S=${OPTARG};;
      # w : window name
        w) W=${OPTARG};;
      # * : handle unknown flags
        *) echo -e "${RED} unknown flag ${flag} ${NC}";;
    esac
done

K=$(grep 'Pubkey:' smckey.pub  | sed 's/🔑 Pubkey: //')
V=$(cat roster.txt)

echo -e "${GREEN}[PUBLISH]${NC} the roster V=${V} on the blockchain using K=${K}"
chaincli --config /tmp/chain1 pool add\
 --key private.key\
 --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Value\
 --args value:key --args \"${K:0:8}\"\
 --args value:value --args \"${K}:${V}\"\
 --args value:command --args WRITE
