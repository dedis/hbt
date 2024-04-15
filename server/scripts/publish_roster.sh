#!/usr/bin/env bash

# This script publishes the SMC roster on the blockchain.

# Requirements:
# from server/blockchain/chaincli: go install

set -e

GREEN='\033[0;32m'    # green color
NC='\033[0m'          # no Color

#K=$(grep 'Pubkey:' smckey.pub  | sed 's/🔑 Pubkey: //' | xxd -r -p | base64)
K=$(grep 'Pubkey:' smckey.pub  | sed 's/🔑 Pubkey: //')
V=$(cat roster.txt)

echo -e "${GREEN}[PUBLISH]${NC} the roster ${V} on the blockchain using key ${K}"

ARGS="--args go.dedis.ch/dela.ContractArg --args go.dedis.ch/dela.Value \
  --args value:key --args ${K} \
  --args value:value --args ${V} \
  --args value:command --args WRITE"

echo -e ${ARGS}
# shellcheck disable=SC2090
LLVL="debug" chaincli --config /tmp/blockchain1 pool add --key /tmp/priv.key ${ARGS}
