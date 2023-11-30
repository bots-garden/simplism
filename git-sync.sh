#!/bin/bash
branch=$(git branch --show-current)

if [[ "$branch" == "main" ]]; then
    echo "👋 You are on the main branch, there is nothing to do."
else
    echo "🚧 Syncing with main..."
    git fetch
    git rebase origin/main
fi