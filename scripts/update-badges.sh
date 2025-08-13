#!/bin/bash

# Simple script to update README badges
# Usage: ./scripts/update-badges.sh

set -e

# Get repository info from git
REPO_URL=$(git remote get-url origin)
if [[ $REPO_URL == *"github.com"* ]]; then
    REPO_OWNER=$(echo "$REPO_URL" | sed 's/.*github\.com[:/]\([^/]*\)\/\([^.]*\).*/\1/')
    REPO_NAME=$(echo "$REPO_URL" | sed 's/.*github\.com[:/]\([^/]*\)\/\([^.]*\).*/\2/')
else
    echo "Not a GitHub repository or remote not set"
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')

# Get coverage if tests have been run
if [ -f "coverage.out" ]; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
else
    COVERAGE="0"
fi

echo "Repository: $REPO_OWNER/$REPO_NAME"
echo "Go version: $GO_VERSION"
echo "Coverage: $COVERAGE%"

# Update README.md
sed -i "s|https://github.com/yourusername/timeseriesdb|https://github.com/$REPO_OWNER/$REPO_NAME|g" README.md
sed -i "s|go-1\.20\+|go-$GO_VERSION|g" README.md
sed -i "s|coverage-0%25|coverage-$COVERAGE%25|g" README.md

echo "README badges updated successfully!"
echo "Don't forget to commit and push your changes."
