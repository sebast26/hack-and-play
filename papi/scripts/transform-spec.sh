#!/bin/bash
# Transform a single OpenAPI spec based on mapping rules
# Usage: ./transform-spec.sh <service-name>
# Example: ./transform-spec.sh orders

set -e

SERVICE=$1
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

if [ -z "$SERVICE" ]; then
    echo "Usage: $0 <service-name>"
    exit 1
fi

SOURCE_SPEC="$ROOT_DIR/mocks/$SERVICE/openapi.yaml"
MAPPING_FILE="$ROOT_DIR/governance/spec-mappings/${SERVICE}-mapping.yaml"
OUTPUT_DIR="$ROOT_DIR/governance/specs/$SERVICE"
OUTPUT_SPEC="$OUTPUT_DIR/openapi.yaml"

# Check files exist
if [ ! -f "$SOURCE_SPEC" ]; then
    echo "Error: Source spec not found: $SOURCE_SPEC"
    exit 1
fi

if [ ! -f "$MAPPING_FILE" ]; then
    echo "Error: Mapping file not found: $MAPPING_FILE"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Read mapping values using yq
SOURCE_BASE_PATH=$(yq '.source_base_path' "$MAPPING_FILE")
TARGET_BASE_PATH=$(yq '.target_base_path' "$MAPPING_FILE")

echo "Transforming $SERVICE spec:"
echo "  Source: $SOURCE_SPEC"
echo "  Mapping: $SOURCE_BASE_PATH -> $TARGET_BASE_PATH"
echo "  Output: $OUTPUT_SPEC"

# Transform paths in the spec
# 1. Replace source_base_path with target_base_path in path keys
# 2. Keep only paths that start with source_base_path

yq eval "
  .paths = (.paths | to_entries | map(
    select(.key | test(\"^${SOURCE_BASE_PATH}\")) |
    .key = (.key | sub(\"^${SOURCE_BASE_PATH}\"; \"${TARGET_BASE_PATH}\"))
  ) | from_entries)
" "$SOURCE_SPEC" > "$OUTPUT_SPEC"

echo "  Done!"
