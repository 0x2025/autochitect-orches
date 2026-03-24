#!/bin/bash
# Run tests for both Node.js and Go
set -e
cd server && npm test && cd ..
cd landlord && go test ./... && cd ..