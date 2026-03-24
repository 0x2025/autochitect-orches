#!/bin/bash
# Build both Node.js and Go components
set -e


cd server && npm run build && cd ..


cd landlord && go build -o ../landlord/bin/landlord ./cmd/main.go && cd ..


