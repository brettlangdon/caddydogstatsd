caddy dogstatsd
===============

A [Caddy](https://caddyserver.com/) plugin for reporting metrics to [Datadog](https://datadoghq.com).

## Installation

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

## Metrics

caddy.response.count - counter - number of requests handled
caddy.response.time - histogram - milliseconds spent handling request
