#!/usr/bin/env bash
set -euo pipefail

retired_devices_file="retired-device-ids.txt"

if [[ ! -s "$retired_devices_file" ]]; then
  exit 0
fi

access_token="$({
  curl --fail-with-body --silent --show-error \
    --data-urlencode "client_id=${TAILSCALE_OAUTH_CLIENT_ID}" \
    --data-urlencode "client_secret=${TAILSCALE_OAUTH_CLIENT_SECRET}" \
    https://api.tailscale.com/api/v2/oauth/token
} | jq --raw-output '.access_token')"

if [[ -z "$access_token" || "$access_token" == "null" ]]; then
  echo "Tailscale OAuth token response did not contain an access token" >&2
  exit 1
fi

while IFS= read -r device_id; do
  [[ -z "$device_id" || "$device_id" == \#* ]] && continue

  status="$(curl --silent --show-error --output /tmp/tailscale-delete-response \
    --write-out '%{http_code}' \
    --request DELETE \
    --header "Authorization: Bearer ${access_token}" \
    "https://api.tailscale.com/api/v2/device/${device_id}")"

  case "$status" in
    2*) echo "Retired Tailscale device ${device_id}" ;;
    404) echo "Tailscale device ${device_id} was already absent" ;;
    *)
      cat /tmp/tailscale-delete-response >&2
      echo "Failed to retire Tailscale device ${device_id}: HTTP ${status}" >&2
      exit 1
      ;;
  esac
done < "$retired_devices_file"
