#!/bin/bash
# Setup script for Autochitect project
# Installs server and landlord dependencies

cd server && npm install

cd ../landlord && go mod tidy

# Setup complete.
