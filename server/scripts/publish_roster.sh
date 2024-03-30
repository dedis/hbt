#!/usr/bin/env bash

# This script publishes the SMC roster on the blockchain.

# Requirements:
# from server/blockchain/chaincli: go install

set -e

GREEN='\033[0;32m'    # green color
NC='\033[0m'          # no Color
KEYFILE=/tmp/priv.key # private key file

K=$(grep 'Pubkey:' smckey.pub  | sed 's/🔑 Pubkey: //' | xxd -r -p | base64)
V=$(cat roster.txt)


echo -e "${GREEN}[PUBLISH]${NC} the roster V=${V} on the blockchain using K=${K}"
chaincli --config /tmp/blockchain1 pool add\
 --key ${KEYFILE}\
 --args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Value\
 --args value:key --args \"${K:0:8}\"\
 --args value:value --args \"${K}:${V}\"\
 --args value:command --args WRITE
