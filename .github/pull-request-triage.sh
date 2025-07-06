#!/usr/bin/env bash

set -e

echo "PR Title: $PR_TITLE"

# Pattern: type(scope?)!?: description
# Examples:
#   feat: add new login feature
#   fix(account)!: breaking change in account handling
#   chore(ci): update workflow
REGEX="^(build|chore|ci|docs|feat|fix|perf|refactor|style|test)(\([^)]+\))?(!)?: .+"

if [[ "$PR_TITLE" =~ $REGEX ]]; then
    echo "✅ PR title follows Conventional Commit format."
else
    # Leave a comment on the PR
    gh pr comment "$PR_NUMBER" --repo "$REPO" --body ":warning: The title of this PR does not follow the [Conventional Commit](https://www.conventionalcommits.org/) format.  

Expected format: \`type(scope?): description\`, e.g. \`feat(login): add new login page\`"

    exit 1
fi

TYPE=$(echo "$PR_TITLE" | sed -E 's/^([a-z]+)(\([^)]+\))?(!)?:.*/\1/')
BREAKING=$(echo "$PR_TITLE" | grep -q '!' && echo "true" || echo "false")

echo "Detected type: $TYPE"
echo "Breaking change: $BREAKING"

if [[ "$BREAKING" == "true" ]]; then
    echo "🛑 Detected breaking change. Adding label."
    gh pr edit "$PR_NUMBER" --repo "$REPO" --add-label "breaking-change"
    exit 0
fi

if [[ "$TYPE" == "feat" ]]; then
    echo "🏷️ Added label: enhancement"
    gh pr edit "$PR_NUMBER" --repo "$REPO" --add-label "enhancement"
    exit 0
fi

if [[ "$TYPE" == "fix" ]]; then
    echo "🏷️ Added label: bug"
    gh pr edit "$PR_NUMBER" --repo "$REPO" --add-label "bug"
    exit 0
fi
