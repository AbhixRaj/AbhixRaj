#!/bin/bash

# Create .devcontainer directory
mkdir -p .devcontainer

# Get the repository name from the current directory
REPO_NAME=$(basename "$(git rev-parse --show-toplevel)")

# Create the devcontainer.json file
cat <<EOL > .devcontainer/devcontainer.json
{
    "name": "$REPO_NAME Codespace",
    "image": "mcr.microsoft.com/vscode/devcontainers/python:3.8",
    "postCreateCommand": "pip install pymongo python-telegram-bot pyTelegramBotAPI certifi && chmod +x /workspaces/$REPO_NAME/*"
}
EOL

# Stage, commit, and push the changes
git add .devcontainer/devcontainer.json
git commit -m "Add devcontainer configuration for $REPO_NAME"
git push origin main

echo "DevContainer setup complete!"
