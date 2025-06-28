#!/usr/bin/env sh

set -e

TOKEN_RESPONSE=$(curl -f --silent --show-error -XPOST "$LAKEKEEPER_AUTH_URL" -d "grant_type=client_credentials" -d "client_id=$LAKEKEEPER_CLIENT_ID" -d "client_secret=$LAKEKEEPER_CLIENT_SECRET")

LAKEKEEPER_TOKEN=$(echo $TOKEN_RESPONSE | jq -r .access_token)

# Print the version, since it is useful debugging information.
curl -f -XPOST -H "Authorization: Bearer $LAKEKEEPER_TOKEN" \
    -H "Content-Type: application/json" \
    -vvv --show-error "$LAKEKEEPER_ENDPOINT/management/v1/bootstrap" \
    -d '{"accept-terms-of-use": true,"is-operator": true}'

echo
