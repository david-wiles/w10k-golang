#!/bin/bash

set -eo pipefail

if [ -z "$1" ]; then
  echo "Usage: ./deploy.sh [domain]";
  exit 1
else
  export TF_VAR_domain=$1
fi

# Build linux executables
./build.sh

# Create the VM
cd tf
terraform apply

# Copy executables to VM
cd ..
scp bin/broadcast "w10k-go.$1":broadcast
scp bin/client2client "w10k-go.$1":client2client
