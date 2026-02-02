#!/bin/bash
# Validate OpenAPI specs with Spectral
# Usage: ./validate.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
RULESET="$ROOT_DIR/governance/spectral/.spectral.yaml"

echo "========================================"
echo "Partner API Spec Validation (Spectral)"
echo "========================================"
echo ""

# Check if spectral is available
if ! npx spectral --version &> /dev/null; then
    echo "Installing Spectral CLI..."
    cd "$ROOT_DIR" && npm install
fi

# Check if ruleset exists
if [ ! -f "$RULESET" ]; then
    echo "Error: Spectral ruleset not found: $RULESET"
    exit 1
fi

ERRORS=0

# Validate individual transformed specs
echo "Validating transformed specs..."
for spec in "$ROOT_DIR"/governance/specs/*/openapi.yaml; do
    if [ -f "$spec" ]; then
        echo "  Linting: ${spec#$ROOT_DIR/}"
        if ! npx spectral lint "$spec" --ruleset "$RULESET"; then
            ERRORS=$((ERRORS + 1))
        fi
        echo ""
    fi
done

# Validate merged Partner API spec
MERGED_SPEC="$ROOT_DIR/governance/output/partner-api.yaml"
if [ -f "$MERGED_SPEC" ]; then
    echo "Validating merged Partner API spec..."
    echo "  Linting: ${MERGED_SPEC#$ROOT_DIR/}"
    if ! npx spectral lint "$MERGED_SPEC" --ruleset "$RULESET"; then
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "Warning: Merged spec not found: $MERGED_SPEC"
    echo "Run ./scripts/aggregate.sh first to generate specs."
fi

echo ""
echo "========================================"
if [ $ERRORS -eq 0 ]; then
    echo "Validation complete! All specs passed."
    exit 0
else
    echo "Validation failed with errors in $ERRORS spec(s)."
    exit 1
fi
