#!/usr/bin/env bash

# This script checks the port used in the system.

sudo netstat -antup | grep 127.0.0.1
