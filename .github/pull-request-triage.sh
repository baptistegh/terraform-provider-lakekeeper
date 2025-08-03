#!/usr/bin/env bash

set -e

echo "PR Title: $PR_TITLE"

TYPE=$(echo "$PR_TITLE" | sed -E 's/^([a-z]+)(\([^)]+\))?(!)?:.*/\1/')
BREAKING=$(echo "$PR_TITLE" | grep -q '!' && echo "true" || echo "false")

echo "Detected type: $TYPE"
echo "Breaking change: $BREAKING"

if [[ "$BREAKING" == "true" ]]; then
    echo "🛑 Detected breaking change. Adding label."
    gh pr edit "$PR_NUMBER" --add-label "breaking-change"
fi

echo "breaking_change=$BREAKING" >> "$GITHUB_OUTPUT"

exit 0