package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const PluginName = "error-status-propagator"

func init() {
	fmt.Printf("[plugin] %s: loaded\n", PluginName)
}

// HandlerRegisterer is the symbol the KrakenD plugin loader looks for.
var HandlerRegisterer = registerer(PluginName)

type registerer string

// RegisterHandlers satisfies the plugin/http-server interface.
func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.newHandler)
}

func (r registerer) newHandler(_ context.Context, extra map[string]interface{}, next http.Handler) (http.Handler, error) {
	// Default prefix matches keys injected by KrakenD's return_error_details feature.
	// E.g. a backend configured with return_error_details: "inventory_check" produces
	// the key "error_inventory_check" in the aggregated response body.
	errorKeyPrefix := "error_"

	// paths is an optional allowlist of URL paths where propagation is active.
	// Empty slice means "apply to all paths" (opt-out model).
	// Non-empty slice means "apply only to listed paths" (opt-in model).
	var paths []string

	if cfg, ok := extra[PluginName].(map[string]interface{}); ok {
		if prefix, ok := cfg["error_key_prefix"].(string); ok && prefix != "" {
			errorKeyPrefix = prefix
		}
		if rawPaths, ok := cfg["paths"].([]interface{}); ok {
			for _, p := range rawPaths {
				if s, ok := p.(string); ok && s != "" {
					paths = append(paths, s)
				}
			}
		}
	}

	if len(paths) == 0 {
		fmt.Printf("[plugin] %s: active on ALL paths, error_key_prefix=%q\n", PluginName, errorKeyPrefix)
	} else {
		fmt.Printf("[plugin] %s: active on %v, error_key_prefix=%q\n", PluginName, paths, errorKeyPrefix)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// If an explicit paths allowlist is configured, skip endpoints not listed.
		if len(paths) > 0 && !pathAllowed(req.URL.Path, paths) {
			next.ServeHTTP(w, req)
			return
		}

		// Capture the response before it is flushed to the client.
		cap := &capture{ResponseWriter: w, buf: &bytes.Buffer{}, code: http.StatusOK}
		next.ServeHTTP(cap, req)

		finalCode := cap.code

		// Only inspect 2xx JSON responses — anything else passes through unchanged.
		if finalCode/100 == 2 && isJSON(cap.Header().Get("Content-Type")) {
			if maxErr := highestBackendError(cap.buf.Bytes(), errorKeyPrefix); maxErr >= 400 {
				finalCode = maxErr
			}
		}

		w.WriteHeader(finalCode)
		w.Write(cap.buf.Bytes()) //nolint:errcheck
	}), nil
}

// pathAllowed reports whether path matches any entry in the allowlist.
// Each entry is matched as an exact path or as a prefix ending with /*.
// Examples:
//
//	"/partner/order-inventory"  — exact match
//	"/partner/*"                — all paths starting with /partner/
func pathAllowed(path string, allowlist []string) bool {
	for _, p := range allowlist {
		if strings.HasSuffix(p, "/*") {
			if strings.HasPrefix(path, strings.TrimSuffix(p, "*")) {
				return true
			}
		} else if path == p {
			return true
		}
	}
	return false
}

// highestBackendError scans the aggregated JSON body for keys with the given
// prefix (injected by KrakenD's return_error_details) and returns the highest
// http_status_code found among them, or 0 if none are present.
func highestBackendError(body []byte, prefix string) int {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return 0
	}

	highest := 0
	for key, val := range data {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		errObj, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		if statusCode, ok := errObj["http_status_code"].(float64); ok {
			if code := int(statusCode); code > highest {
				highest = code
			}
		}
	}
	return highest
}

func isJSON(contentType string) bool {
	return strings.Contains(contentType, "application/json")
}

// capture wraps http.ResponseWriter to buffer the response so the status code
// and body can be inspected and modified before they reach the client.
type capture struct {
	http.ResponseWriter
	buf  *bytes.Buffer
	code int
}

func (c *capture) WriteHeader(code int) { c.code = code }
func (c *capture) Write(b []byte) (int, error) { return c.buf.Write(b) }
