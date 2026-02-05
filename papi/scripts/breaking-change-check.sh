#!/bin/bash
# Check for breaking changes between baseline and current Partner API spec
# Usage: ./breaking-change-check.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BASELINE_SPEC="$ROOT_DIR/governance/baseline/partner-api.yaml"
CURRENT_SPEC="$ROOT_DIR/governance/output/partner-api.yaml"

echo "========================================"
echo "Partner API Breaking Change Detection"
echo "========================================"
echo ""

# Function to install oasdiff
install_oasdiff() {
    echo "Installing oasdiff..."

    # Detect OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "Error: Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    # macOS uses universal binary (darwin_all)
    if [ "$OS" = "darwin" ]; then
        ARCH="all"
    fi

    # Get latest version (use -L to follow redirects)
    VERSION=$(curl -sL https://api.github.com/repos/oasdiff/oasdiff/releases/latest | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest oasdiff version"
        exit 1
    fi

    DOWNLOAD_URL="https://github.com/oasdiff/oasdiff/releases/download/v${VERSION}/oasdiff_${VERSION}_${OS}_${ARCH}.tar.gz"

    echo "  Downloading oasdiff v${VERSION} for ${OS}/${ARCH}..."

    # Create bin directory if needed
    mkdir -p "$ROOT_DIR/bin"

    # Download and extract
    curl -sL "$DOWNLOAD_URL" | tar -xz -C "$ROOT_DIR/bin" oasdiff
    chmod +x "$ROOT_DIR/bin/oasdiff"

    echo "  Installed to $ROOT_DIR/bin/oasdiff"
}

# Check if oasdiff is available
OASDIFF=""
if command -v oasdiff &> /dev/null; then
    OASDIFF="oasdiff"
elif [ -x "$ROOT_DIR/bin/oasdiff" ]; then
    OASDIFF="$ROOT_DIR/bin/oasdiff"
else
    install_oasdiff
    OASDIFF="$ROOT_DIR/bin/oasdiff"
fi

echo "Using: $($OASDIFF --version 2>&1 | head -1)"
echo ""

# Check if baseline spec exists
if [ ! -f "$BASELINE_SPEC" ]; then
    echo "Error: Baseline spec not found: $BASELINE_SPEC"
    echo ""
    echo "To create a baseline, run:"
    echo "  cp governance/output/partner-api.yaml governance/baseline/partner-api.yaml"
    exit 1
fi

# Check if current spec exists
if [ ! -f "$CURRENT_SPEC" ]; then
    echo "Error: Current spec not found: $CURRENT_SPEC"
    echo "Run ./scripts/aggregate.sh first to generate the spec."
    exit 1
fi

echo "Comparing specs:"
echo "  Baseline: ${BASELINE_SPEC#$ROOT_DIR/}"
echo "  Current:  ${CURRENT_SPEC#$ROOT_DIR/}"
echo ""

# Run oasdiff breaking change detection
echo "Checking for breaking changes..."
echo ""

# oasdiff breaking returns exit code 1 if breaking changes found
# --fail-on WARN treats warnings (like removing optional fields) as breaking changes
if $OASDIFF breaking "$BASELINE_SPEC" "$CURRENT_SPEC" --format text --fail-on WARN; then
    echo "========================================"
    echo "No breaking changes detected."
    echo "========================================"
    exit 0
else
    EXIT_CODE=$?
    echo ""
    echo "========================================"
    echo "BREAKING CHANGES DETECTED!"
    echo "========================================"
    echo ""
    echo "Next steps:"
    echo "  1. If this is intentional, coordinate with the platform team"
    echo "  2. Post in #partner-api-changes Slack channel"
    echo "  3. Create a JIRA ticket with change details"
    echo "  4. See: docs/breaking-change-workflow.md"
    echo ""
    exit $EXIT_CODE
fi
