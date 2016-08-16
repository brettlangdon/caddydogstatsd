caddy dogstatsd
===============

A [Caddy](https://caddyserver.com/) middleware plugin for reporting metrics to [Datadog](https://datadoghq.com).

## Installation

**Coming soon**

## Configuration

```
dogstatsd [{host:port} [samplerate]]
```

```
dogstatsd {
  host {host:port}
  samplerate {samplerate}
  namespace {namespace}
  tags {name:value} [{name:value}...]
}
```

### Configuration options

- `host` - the `host:port` where `dogstatsd` metrics should be sent - default: `127.0.0.1:8125`
- `samplerate` - a `float` indicating the sample rate of requests to record metrics for
  - For example, a `samplerate` of `0.5` means metrics will be emitted for only half of the requests
- `namespace` - an optional namespace to prepend to each metric emitted - by default there is none
  - For example, a `namespace` of `my_app` will yield `my_app.caddy.response.time` as a metric
- `tags` - an optional list of global tags to set for this server - by default there are none
  - Example tags, `env:production`, `app:my_app`

### Example config
```
dogstatsd {
  samplerate 0.75
  namespace my_app
  tags env:production service:caddy
}
```

## Metrics

- `[namespace.]caddy.response.count` - counter - number of requests handled
- `[namespace.]caddy.response.time` - histogram - milliseconds spent handling request
