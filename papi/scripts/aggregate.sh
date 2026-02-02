#!/bin/bash
# Spec Aggregation Pipeline
# Transforms, merges OpenAPI specs, and generates KrakenD config

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "========================================"
echo "Partner API Spec Aggregation Pipeline"
echo "========================================"
echo ""

# Check dependencies
check_dependency() {
    if ! command -v "$1" &> /dev/null; then
        echo "Error: $1 is required but not installed."
        echo "Install with: $2"
        exit 1
    fi
}

check_dependency "yq" "brew install yq"
check_dependency "jq" "brew install jq"

# Check if openapi-merge-cli is available
if ! npx openapi-merge-cli --version &> /dev/null; then
    echo "Installing openapi-merge-cli..."
    cd "$ROOT_DIR" && npm install
fi

echo "Step 1: Transform specs from mocks to governance/specs"
echo "-------------------------------------------------------"

# Get list of services from manifest
SERVICES=$(yq '.sources[].name' "$ROOT_DIR/governance/manifest.yaml")

for service in $SERVICES; do
    enabled=$(yq ".sources[] | select(.name == \"$service\") | .enabled" "$ROOT_DIR/governance/manifest.yaml")
    if [ "$enabled" = "true" ]; then
        "$SCRIPT_DIR/transform-spec.sh" "$service"
    else
        echo "Skipping disabled service: $service"
    fi
done

echo ""
echo "Step 2: Merge specs into partner-api.yaml"
echo "------------------------------------------"

cd "$ROOT_DIR/governance"
npx openapi-merge-cli --config openapi-merge.yaml

echo "  Output: governance/output/partner-api.yaml"

echo ""
echo "Step 3: Generate KrakenD configuration"
echo "--------------------------------------"

"$SCRIPT_DIR/generate-krakend.sh"

echo ""
echo "========================================"
echo "Pipeline complete!"
echo "========================================"
echo ""
echo "Generated files:"
echo "  - governance/specs/*/openapi.yaml  (transformed specs)"
echo "  - governance/output/partner-api.yaml (merged Partner API spec)"
echo "  - gateway/krakend.json (KrakenD configuration)"
echo ""
echo "To test: docker-compose down && docker-compose up -d"
