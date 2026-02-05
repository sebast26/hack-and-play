#!/bin/bash
# Run Schemathesis API testing against a running service
# Usage: ./schemathesis-test.sh [options]
#
# Options:
#   --spec PATH/URL      Path or URL to OpenAPI spec (required)
#   --url URL            Base URL of the API to test (required)
#   --timeout MINUTES    Test timeout in minutes (default: 60)
#   --retries N          Number of retry attempts (default: 2)
#   --output DIR         Output directory for reports (default: ./schemathesis-reports)
#   --checks CHECKS      Comma-separated checks to run (default: core checks)
#   --dry-run            Show what would be tested without running

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Default values
SPEC_PATH=""
BASE_URL=""
TIMEOUT_MINUTES=60
RETRIES=2
OUTPUT_DIR="$ROOT_DIR/schemathesis-reports"
DRY_RUN=false
# Default checks: exclude 'unsupported_method' which tests TRACE and other exotic methods
CHECKS="not_a_server_error,status_code_conformance,content_type_conformance,response_headers_conformance,response_schema_conformance"

# Schemathesis Docker image
SCHEMATHESIS_IMAGE="schemathesis/schemathesis:stable"

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --spec)
            SPEC_PATH="$2"
            shift 2
            ;;
        --url)
            BASE_URL="$2"
            shift 2
            ;;
        --timeout)
            TIMEOUT_MINUTES="$2"
            shift 2
            ;;
        --retries)
            RETRIES="$2"
            shift 2
            ;;
        --output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        --checks)
            CHECKS="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  --spec PATH/URL      Path or URL to OpenAPI spec (required)"
            echo "  --url URL            Base URL of the API to test (required)"
            echo "  --timeout MINUTES    Test timeout in minutes (default: 60)"
            echo "  --retries N          Number of retry attempts (default: 2)"
            echo "  --output DIR         Output directory for reports"
            echo "  --checks CHECKS      Comma-separated checks to run"
            echo "  --dry-run            Show what would be tested without running"
            echo ""
            echo "Available checks:"
            echo "  not_a_server_error, status_code_conformance, content_type_conformance,"
            echo "  response_headers_conformance, response_schema_conformance,"
            echo "  negative_data_rejection, positive_data_acceptance, unsupported_method,"
            echo "  use_after_free, ensure_resource_availability, ignored_auth, all"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate required arguments
if [ -z "$SPEC_PATH" ]; then
    echo "Error: --spec is required"
    exit 1
fi

if [ -z "$BASE_URL" ]; then
    echo "Error: --url is required"
    exit 1
fi

echo "========================================"
echo "Partner API Implementation Testing"
echo "========================================"
echo ""
echo "Spec:        $SPEC_PATH"
echo "Base URL:    $BASE_URL"
echo "Timeout:     ${TIMEOUT_MINUTES} minutes"
echo "Retries:     $RETRIES"
echo "Output:      $OUTPUT_DIR"
echo "Checks:      $CHECKS"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is required but not installed."
    exit 1
fi

# Pull Schemathesis image if needed
if ! docker image inspect "$SCHEMATHESIS_IMAGE" &> /dev/null; then
    echo "Pulling Schemathesis Docker image..."
    docker pull "$SCHEMATHESIS_IMAGE"
fi

# Build Schemathesis command
build_schemathesis_cmd() {
    local cmd="docker run --rm"

    # Network mode for accessing host services
    cmd="$cmd --network host"

    # Mount output directory
    cmd="$cmd -v $OUTPUT_DIR:/output"

    # Mount spec file if it's a local path
    if [ -f "$SPEC_PATH" ]; then
        local abs_path
        abs_path=$(realpath "$SPEC_PATH")
        cmd="$cmd -v $abs_path:/spec.yaml:ro"
        cmd="$cmd $SCHEMATHESIS_IMAGE run /spec.yaml"
    else
        cmd="$cmd $SCHEMATHESIS_IMAGE run \"$SPEC_PATH\""
    fi

    # Base URL
    cmd="$cmd --url \"$BASE_URL\""

    # Checks to run
    cmd="$cmd --checks $CHECKS"

    # Output options
    cmd="$cmd --report junit --report-dir /output"

    # Request timeout (in seconds)
    cmd="$cmd --request-timeout=10"

    # Limit max failures to avoid long runs
    cmd="$cmd --max-failures=20"

    echo "$cmd"
}

# Function to run tests with retries
run_tests() {
    local attempt=1
    local max_attempts=$((RETRIES + 1))
    local timeout_seconds=$((TIMEOUT_MINUTES * 60))

    while [ $attempt -le $max_attempts ]; do
        echo "Test attempt $attempt of $max_attempts..."
        echo ""

        local cmd
        cmd=$(build_schemathesis_cmd)

        if [ "$DRY_RUN" = true ]; then
            echo "Would run:"
            echo "$cmd"
            return 0
        fi

        # Run with timeout
        set +e
        timeout "$timeout_seconds" bash -c "$cmd"
        local exit_code=$?
        set -e

        if [ $exit_code -eq 0 ]; then
            echo ""
            echo "========================================"
            echo "All tests passed!"
            echo "========================================"
            echo ""
            echo "Reports saved to: $OUTPUT_DIR"
            return 0
        elif [ $exit_code -eq 124 ]; then
            echo ""
            echo "Test run timed out after ${TIMEOUT_MINUTES} minutes"
        else
            echo ""
            echo "Tests failed with exit code: $exit_code"
        fi

        if [ $attempt -lt $max_attempts ]; then
            echo "Retrying in 10 seconds..."
            sleep 10
        fi

        attempt=$((attempt + 1))
    done

    echo ""
    echo "========================================"
    echo "Tests failed after $max_attempts attempts"
    echo "========================================"
    echo ""
    echo "Check reports in: $OUTPUT_DIR"
    return 1
}

# Wait for service to be ready
wait_for_service() {
    local url="$1"
    local max_wait=60
    local waited=0

    echo "Waiting for service at $url..."

    # Try common health check endpoints
    local endpoints=("/__health" "/health" "/partner/orders" "/partner/users" "/")

    while [ $waited -lt $max_wait ]; do
        for endpoint in "${endpoints[@]}"; do
            if curl -sf "${url}${endpoint}" > /dev/null 2>&1; then
                echo "Service is ready! (checked ${endpoint})"
                return 0
            fi
        done
        sleep 2
        waited=$((waited + 2))
    done

    echo "Service did not become ready within ${max_wait} seconds"
    return 1
}

# Main execution
echo "Checking service availability..."
if ! wait_for_service "$BASE_URL"; then
    echo "Error: Service is not available at $BASE_URL"
    exit 1
fi

echo ""
if [ -f "$SPEC_PATH" ]; then
    echo "Using local spec: $SPEC_PATH"
else
    echo "Fetching spec from $SPEC_PATH..."
    if ! curl -sf "$SPEC_PATH" > /dev/null 2>&1; then
        echo "Error: Could not fetch spec from $SPEC_PATH"
        exit 1
    fi
    echo "Spec is accessible!"
fi

echo ""
echo "Starting Schemathesis tests..."
echo ""

run_tests
