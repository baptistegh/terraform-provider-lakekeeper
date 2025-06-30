#!/usr/bin/env sh

CONTAINER_ENGINE="${CONTAINER_ENGINE:-docker}"

set -e

if [ "$CONTAINER_ENGINE" != "docker" ]; then
  echo "Using container engine $CONTAINER_ENGINE"
fi

printf 'Waiting for Lakekeeper container to become healthy'

until test -n "$($CONTAINER_ENGINE ps --quiet --filter label=terraform-provider-lakekeeper/owned --filter health=healthy)"; do
  printf '.'
  sleep 5
done

echo
echo "Lakekeeper is healthy at $LAKEKEEPER_ENDPOINT"

# Get token
echo "Getting OIDC access token for Lakekeeper"
TOKEN=$(curl --silent --show-error --fail \
  --data "scope=lakekeeper&grant_type=client_credentials&client_id=$LAKEKEEPER_CLIENT_ID&client_secret=$LAKEKEEPER_CLIENT_SECRET" \
  "$LAKEKEEPER_AUTH_URL" | jq -r '.access_token')

# Print the server info, since it is useful debugging information.
echo "Lakekeeper server info:"
curl --fail --show-error --silent -H "Authorization: Bearer $TOKEN" "$LAKEKEEPER_ENDPOINT/management/v1/info"
echo
