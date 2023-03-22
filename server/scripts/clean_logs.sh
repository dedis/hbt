#!/usr/bin/env bash

# This script cleans all local logs in linux.

set -e

GREEN='\033[0;32m'    # green color
RED='\033[1;31;46m'   # red color
NC='\033[0m'          # no Color

echo -e "${RED} CLEANING ALL LOGS IN CURRENT DIRECTORY - READY ? ${NC}"
read -p "Press Enter to continue" </dev/tty

rm *.log

