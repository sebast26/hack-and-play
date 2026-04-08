# Error Status Propagation in Multi-Backend Aggregation

## Problem

KrakenD aggregates responses from multiple backends into a single endpoint. When one backend returns an error (4xx/5xx), KrakenD's default behavior does **not** propagate that backend's HTTP status to the client.

| Scenario | Default KrakenD behavior |
|---|---|
| All backends → 200 | Client receives 200 |
| Some backends → error, some → 200 | Client receives **200** (partial response) |
| All backends → error | Client receives 500 |

This means a backend returning 400 is silently "swallowed" — the client has no way to detect the error from the HTTP status alone.

### What `return_error_details` Gives You

KrakenD CE provides `return_error_details` at the backend level. When configured, it includes the backend's error in the aggregated response body under a key like `error_<alias>`:

```json
{
  "orders": { "orders": [...] },
  "error_inventory_check": {
    "http_status_code": 400,
    "http_body": "{\"error\":\"inventory_unavailable\",...}"
  }
}
```

The body exposes the error, but **the HTTP status remains 200**. Partners cannot rely on HTTP status codes to detect backend failures without parsing the body.

### Why `return_error_code` Doesn't Help Here

KrakenD CE provides `return_error_code: true` at the backend level, but this only works for **single-backend endpoints**. In a multi-backend aggregated endpoint it has no effect.

---

## Solution: Custom `http-server` Plugin

We implemented a KrakenD Go plugin that wraps every response before it reaches the client. For enabled endpoints, it inspects the aggregated JSON body for `error_*` keys (added by `return_error_details`) and promotes the highest backend error code to the response's HTTP status.

### How It Works

```
Request
  └─► KrakenD aggregates backends
        ├─ orders service      → HTTP 200  ✓
        └─ inventory service   → HTTP 400  ✗
            (body includes error_inventory_check)
  └─► Plugin intercepts response (status 200, body with error key)
        ├─ Finds error_inventory_check.http_status_code = 400
        └─ Rewrites response status to 400
  └─► Client receives HTTP 400
```

The plugin:
1. Buffers the full response before sending it to the client
2. If the response is 2xx JSON, scans top-level keys for the configured prefix (`error_` by default)
3. If any error key contains `http_status_code ≥ 400`, uses the highest code as the final HTTP status
4. Sends the body unchanged — only the status is modified

Endpoints without `return_error_details` or without errors produce no `error_*` keys, so the plugin passes them through unaffected.

### Plugin Configuration

The plugin is loaded at the **service level** but activates only on **explicitly listed paths** (opt-in):

```json
{
  "plugin": {
    "pattern": ".so",
    "folder": "/opt/krakend/plugins/"
  },
  "extra_config": {
    "plugin/http-server": {
      "name": ["error-status-propagator"],
      "error-status-propagator": {
        "error_key_prefix": "error_",
        "paths": [
          "/partner/order-inventory"
        ]
      }
    }
  }
}
```

| Config key | Default | Description |
|---|---|---|
| `error_key_prefix` | `"error_"` | Prefix of keys to inspect in the response body |
| `paths` | *(absent)* | Opt-in allowlist; if omitted, applies to all endpoints |

**Path patterns:**

| Entry | Matches |
|---|---|
| `"/partner/order-inventory"` | Exact path only |
| `"/partner/*"` | All paths starting with `/partner/` |
| *(field absent)* | All endpoints |

### Enabling on a New Endpoint

Two changes required — no plugin rebuild needed:

**1. Add `return_error_details` to the failing backend:**

```json
{
  "endpoint": "/partner/checkout",
  "backend": [
    {
      "url_pattern": "/api/v1/cart",
      "host": ["http://cart:8084"],
      "encoding": "json",
      "group": "cart"
    },
    {
      "url_pattern": "/api/v1/payment/validate",
      "host": ["http://payments:8085"],
      "encoding": "json",
      "group": "payment",
      "extra_config": {
        "backend/http": {
          "return_error_details": "payment_check"
        }
      }
    }
  ]
}
```

**2. Add the path to the plugin's allowlist:**

```json
"paths": [
  "/partner/order-inventory",
  "/partner/checkout"
]
```

---

## Build Considerations

Building a KrakenD plugin requires the **exact same Go version** that KrakenD itself was compiled with. Go's plugin loader verifies a runtime package hash at load time and rejects mismatches with:

```
plugin was built with a different version of package runtime
```

### Version Chain for KrakenD 2.5

| Component | Version |
|---|---|
| KrakenD 2.5 binary | Go 1.20.13, Alpine 3.18, musl |
| Alpine 3.18 `go` package | Go 1.20.11 ← **wrong version** |
| Official go.dev binary | Go 1.20.13, glibc ← **wrong libc** |

Neither standard source works out of the box. Our Dockerfile handles this:

1. **Starts from `devopsfaith/krakend:2.5`** — ensures identical Alpine 3.18 / musl environment
2. **Installs `gcompat`** — provides glibc→musl compatibility shim so the glibc-linked Go 1.20.13 binary from go.dev can execute on Alpine
3. **Downloads Go 1.20.13 from go.dev** — exact version match for the runtime hash check
4. **Uses `CC="gcc -fuse-ld=gold"`** — works around an arm64 Alpine issue where `collect2` cannot locate the default `ld.bfd` linker; switching to `ld.gold` (provided by `binutils-gold`) resolves it

```dockerfile
FROM devopsfaith/krakend:2.5 AS plugin-builder

RUN apk add --no-cache musl-dev gcc binutils binutils-gold wget ca-certificates gcompat && \
    wget -q "https://go.dev/dl/go1.20.13.linux-${TARGETARCH}.tar.gz" -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go

WORKDIR /plugin
COPY plugins/error-status-propagator/ .
RUN go version && CGO_ENABLED=1 CC="gcc -fuse-ld=gold" \
    go build -buildmode=plugin -o /tmp/error-status-propagator.so .

FROM devopsfaith/krakend:2.5
COPY --from=plugin-builder /tmp/error-status-propagator.so /opt/krakend/plugins/
```

### Upgrading KrakenD

When upgrading KrakenD, update `ARG GO_VERSION` in the Dockerfile to match the new binary's Go version:

```sh
docker run --rm devopsfaith/krakend:<new-version> krakend version
# Look for: Go Version: 1.xx.yy
```

Then rebuild the gateway image.
