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

# Print the version, since it is useful debugging information.
curl --silent --show-error "$LAKEKEEPER_ENDPOINT/management/v1/info"
echo
