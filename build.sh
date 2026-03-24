#!/bin/bash
# Build both Node.js and Go components
set -e

echo "Building Node.js server..."
cd server && npm run build && cd ..

echo "Building Go landlord..."
cd landlord && go build -o ../landlord/bin/landlord ./cmd/main.go && cd ..

echo "Build completed."
