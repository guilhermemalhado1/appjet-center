#!/bin/bash

# This script builds the appjet-client-cli executable

# Get today's date and time in the format MMDDYYYYHHMM
datetime=$(date +"%m%d%Y%H%M")

cd appjet-decision-provider

# Create a build folder with the appended date and time
go build -o "../artifacts/appjet-decision-provider-$datetime/appjet-decision-provider" .
#go build -o "../artifacts/appjet-server-daemon-$datetime/appjet-server-daemon" .