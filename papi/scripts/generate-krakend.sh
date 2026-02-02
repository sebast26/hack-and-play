#!/bin/bash
# Generate KrakenD configuration from merged Partner API spec
# Usage: ./generate-krakend.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

PARTNER_API_SPEC="$ROOT_DIR/governance/output/partner-api.yaml"
KRAKEND_CONFIG="$ROOT_DIR/gateway/krakend.json"
MAPPINGS_DIR="$ROOT_DIR/governance/spec-mappings"

if [ ! -f "$PARTNER_API_SPEC" ]; then
    echo "Error: Partner API spec not found: $PARTNER_API_SPEC"
    echo "Run the aggregation pipeline first."
    exit 1
fi

echo "Generating KrakenD configuration..."
echo "  Source: $PARTNER_API_SPEC"
echo "  Output: $KRAKEND_CONFIG"

# Function to find backend info for a given path
find_backend_for_path() {
    local path="$1"
    for mapping_file in "$MAPPINGS_DIR"/*-mapping.yaml; do
        target_path=$(yq '.target_base_path' "$mapping_file")
        if [[ "$path" == "$target_path"* ]]; then
            echo "$mapping_file"
            return 0
        fi
    done
    return 1
}

# Start building KrakenD config
cat > "$KRAKEND_CONFIG" << 'HEADER'
{
  "$schema": "https://www.krakend.io/schema/krakend.json",
  "version": 3,
  "name": "Partner API Gateway",
  "timeout": "5s",
  "cache_ttl": "0s",
  "output_encoding": "json",
  "port": 8080,
  "extra_config": {
    "router": {
      "return_error_msg": true
    }
  },
  "endpoints": [
HEADER

# Extract paths from the merged spec and generate endpoints
FIRST=true
yq -o=json '.paths | to_entries | .[]' "$PARTNER_API_SPEC" | jq -c '{path: .key, methods: (.value | keys)}' | while IFS= read -r endpoint; do
    path=$(echo "$endpoint" | jq -r '.path')
    methods=$(echo "$endpoint" | jq -r '.methods[]')

    # Find matching mapping file
    mapping_file=$(find_backend_for_path "$path") || {
        echo "Warning: No backend mapping found for path: $path" >&2
        continue
    }

    # Read backend info from mapping
    backend_host=$(yq '.backend.host' "$mapping_file")
    source_base=$(yq '.source_base_path' "$mapping_file")
    target_base=$(yq '.target_base_path' "$mapping_file")

    # Convert partner path back to internal path
    internal_path="${path/$target_base/$source_base}"

    for method in $methods; do
        # Skip non-HTTP methods
        case "$method" in
            get|post|put|delete|patch|head|options)
                METHOD_UPPER=$(echo "$method" | tr '[:lower:]' '[:upper:]')

                if [ "$FIRST" = true ]; then
                    FIRST=false
                else
                    echo "," >> "$KRAKEND_CONFIG"
                fi

                cat >> "$KRAKEND_CONFIG" << ENDPOINT
    {
      "endpoint": "$path",
      "method": "$METHOD_UPPER",
      "output_encoding": "json",
      "backend": [
        {
          "url_pattern": "$internal_path",
          "host": ["$backend_host"],
          "encoding": "json"
        }
      ]
    }
ENDPOINT
                ;;
        esac
    done
done

# Close the JSON
cat >> "$KRAKEND_CONFIG" << 'FOOTER'

  ]
}
FOOTER

echo "  Done! Generated $(grep -c '"endpoint"' "$KRAKEND_CONFIG") endpoints."
