#!/usr/bin/env bash

# This script tests the registration process

# send the data to the registration server
curl -F "name=John Doe" -F "passport=12AB456789" -F "role=0" -F "image=@./test.jpg" -F "registered=false" localhost:3000/document

# send the data to the registration server
curl -F "name=John Doe" -F "passport=12AB456789" -F "role=0" -F "image=@./test.jpg" -F "registered=false" localhost:3000/document

