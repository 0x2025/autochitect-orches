#!/bin/bash
# Run tests for both Node.js and Go
set -e

echo "Running Node.js tests..."
cd server && npm test && cd ..

echo "Running Go tests..."
cd landlord && go test ./... && cd ..

echo "All tests passed."
