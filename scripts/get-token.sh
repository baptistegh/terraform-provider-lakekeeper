#!/usr/bin/env sh

set -e 

curl -XPOST --silent --fail --show-error http://localhost:30080/realms/iceberg/protocol/openid-connect/token \
    -d client_id=lakekeeper-admin \
    -d client_secret=KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL \
    -d grant_type=client_credentials \
    -d scope=lakekeeper | jq -r .access_token