# tsgen

`tsgen` is a little Go program to simulate HTTP requests faults and demonstrate how Prometheus alerts based on the [Multiwindow, Multi-Burn-Rate Alerts](https://sre.google/workbook/alerting-on-slos/) works.
Using it is easy, just define some paths, attach some faults and start to make requests.

## Config
`tsgen` requires a YAML config file, to define a list of paths; each path has some faults, formed by a HTTP Status Code and error rate.
Also every path should include a response time, in milliseconds.

e.g.:
```
paths:
  - name: hello
    response_time: 20ms
    faults:
      - code: 500
        rate: 0.001
      - code: 401
        rate: 0.1

  - name: world
    response_time: 10ms
    faults:
      - code: 502
        rate: 0.005
```

## Usage

To build `tsgen` Docker images just run

```
docker build -t tsgen .
```

A docker-compose file is included, so to start `tsgen` with Prometheus and Grafana just execute

```
docker-compose up
```

To login in Grafana use default config values
```
username: admin
password: admin
```

## API

`tsgen` has a very simple API, as follows

- GET /api/health
- GET /api/faults
- GET /api/config
- POST /api/paths/:path/faults/:code

The POST option is designed to changes faults errors rates at runtime, using a JSON payload

```
{"rate": float32}
```

e.g.:
```
curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/paths/hello/faults/500 --data '{"rate": 0.1}'
```

## Examples
Three examples are included in a directory with the same name:
- basic
- errors-rates
- multi-shop

As its name points, *basic* is the most trivial, just to show how tsgen works.

The *errors-rates* is designed to demonstrate how the [Multiwindow, Multi-Burn-Rate Alerts](https://sre.google/workbook/alerting-on-slos/) works, particularly how alerts are triggered in correspondence with the error rate value; lower values require more time while greater ones are triggered very quickly.

The last case, *multi-shop*, is oriented to show how the Grafana Services dashboard is useful for multiple services when some conventions in the metrics names are in place.
