#!/bin/bash
# Setup script for Autochitect project
# Installs server and landlord dependencies

echo "Installing server dependencies..."
cd server && npm install

echo "Installing landlord dependencies..."
cd ../landlord && go mod tidy

echo "Setup complete."
