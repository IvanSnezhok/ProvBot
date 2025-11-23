#!/bin/bash

# --- Configuration ---
# Replace these with your actual server details
SERVER_USER="root"
SERVER_HOST="your_server_ip"
DEPLOY_DIR="/opt/provbot"
# ---------------------

echo "Building for Linux..."
make build-linux

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Deploying to $SERVER_USER@$SERVER_HOST:$DEPLOY_DIR..."

# Create directory if it doesn't exist
ssh $SERVER_USER@$SERVER_HOST "mkdir -p $DEPLOY_DIR"

# Upload binary
scp bin/provbot-linux $SERVER_USER@$SERVER_HOST:$DEPLOY_DIR/

# Upload .env (be careful not to overwrite if you have production secrets there, 
# maybe upload .env.example or manage secrets manually)
# scp .env $SERVER_USER@$SERVER_HOST:$DEPLOY_DIR/

# Upload configs
scp -r configs $SERVER_USER@$SERVER_HOST:$DEPLOY_DIR/

# Upload systemd service file
scp deploy/provbot.service $SERVER_USER@$SERVER_HOST:/etc/systemd/system/

echo "Reloading systemd and restarting service..."
ssh $SERVER_USER@$SERVER_HOST "systemctl daemon-reload && systemctl restart provbot && systemctl status provbot"

echo "Deployment complete!"
