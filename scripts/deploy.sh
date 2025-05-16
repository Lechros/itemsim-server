#!/bin/bash

# Exit on error
set -e

# Configuration
SERVICE_NAME="itemsim-server"
IMAGE_NAME="itemsim-server"

# Log function
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Error handling
handle_error() {
    log "Error occurred in deployment script"
    exit 1
}

trap handle_error ERR

# Load the Docker image and capture the output
log "Loading Docker image..."
docker load < "$IMAGE_PATH"

# Update the service
log "Updating Docker service..."
docker service update "$SERVICE_NAME" --image "$IMAGE_NAME:$IMAGE_TAG"

# Cleanup
log "Cleaning up..."
docker image prune -f
rm "$IMAGE_PATH"

log "Deployment completed successfully"
