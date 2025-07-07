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
    gh pr edit "$PR_NUMBER" --repo "$REPO" --remove-label "invalid-title" 2> /dev/null || echo "error removing label"
else
    # Tag invalid title
    LABEL_EXISTS=$(gh pr view "$PR_NUMBER" --repo "$REPO" --json labels --jq '.labels[].name | select(. == "invalid-title")')

    if [ -z "$LABEL_EXISTS" ]; then
         # Leave a comment on the PR
        gh pr comment "$PR_NUMBER" --repo "$REPO" --body ":warning: The title of this PR does not follow the [Conventional Commit](https://www.conventionalcommits.org/) format.  

Expected format: \`type(scope?): description\`, e.g. \`feat(login): add new login page\`"

        gh pr edit $PR_NUMBER --add-label "invalid-title"
    else
        echo "Label invalid-title present, skipping comment."
    fi

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

if [[ "$TYPE" == "docs" ]]; then
    echo "🏷️ Added label: documentation"
    gh pr edit "$PR_NUMBER" --repo "$REPO" --add-label "documentation"
    exit 0
fi

if [[ "$TYPE" == "ci" ]]; then
    echo "🏷️ Added label: ci"
    gh pr edit "$PR_NUMBER" --repo "$REPO" --add-label "ci"
    exit 0
fi
