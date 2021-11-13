# tsgen

`tsgen` is a little Go program to simulate HTTP requests faults and show how Prometheus alerts based on the [Multiwindow, Multi-Burn-Rate Alerts](https://sre.google/workbook/alerting-on-slos/) works.
Using it is easy, just define some paths, attach some faults and start to make requests.

## Config
`tsgen` requires a YAML config file to define a list of paths; each path has some faults, formed by a HTTP Status Code with an error rate.

e.g.:
```
paths:
  - name: hello
    faults:
      - code: 500
        rate: 0.001
      - code: 401
        rate: 0.1

  - name: world
    faults:
      - code: 502
        rate: 0.005
```

## API

`tsgen` has a very simple API, as follows

- GET `/api/health`
- GET `/api/faults`
- GET `/api/config`
- PUT `/api/paths/:path/faults/:code`

The last option is designed to changes faults errors rates at runtime, using a JSON payload

```
{"rate": float32}
```

e.g.:
```
curl -X PUT -H "Content-Type: application/json" http://localhost:8080/api/paths/hello/faults/500 --data '{"rate": 0.1}'
```

## Usage

To build `tsgen` Docker image just run

```
docker build -t tsgen .
```

## Examples
Three examples are included in a directory with the same name:
- basic
- errors-rates
- multi-shop

As its name points, *basic* is the most trivial, just to show how `tsgen` works.

The *errors-rates* example is designed to prove how the [Multiwindow, Multi-Burn-Rate Alerts](https://sre.google/workbook/alerting-on-slos/) works.

The last case, *multi-shop*, is intended to show how the Grafana "Services" dashboard is useful for multiple services when some conventions in the metrics names and labels are followed.

A directory named `_common` contains common Prometheus rules (`_common/promtheus/etc`). and Grafana Dashboards (`_common/grafana/dashboards`).

### Errors Rates
As the two other examples, a structure similar to the next is included

```
├── docker-compose.yml
├── prometheus
│   ├── data
│   └── etc
│       └── prometheus.yml
└── tsgen
    └── config.yml
```

To run it go to the corresponding directory first
```
cd examples/errors-rates
```
and then execute
```
docker-compose up
```

after that, the following services will start:
- `tsgen` at [localhost:8080](http://localhost:8080)
- Prometheus at [localhost:9090](http://localhost:9090)
- Grafana at [localhost:3000](http://localhost:3000)
- Pushgateway (if included) at [localhost:9091](http://localhost:9091)

To login in Grafana use default config values
```
username: admin
password: admin
```

The `tsgen` config file for this example contains nine paths, following a sequence of Fibonacci numbers (from 2 to 89): `e0002pct`, `e0008pct` until `e0089pct`.

For each path, the error rate is set to 0.001, and during the simulation is changed (using the API) to a value corresponding to the sequence number divided by 100. This way the path `e0002pct` will respond with a 0.02 rate of errors while the last path fails almost all requests. The direct consequence of this is observed in the time required to trigger the corresponding alert for each path.

To start making requests a load generator is useful; in particular a Constant Throughput style is convenient to keep the number of requests stable. For this purpose https://github.com/jrhrmsll/stress is used, but other solutions are available.

With `stress` a script similar to the next could be used
```
#!/usr/bin/env bash

paths=("e0002pct" "e0003pct" "e0005pct" "e0008pct" "e0013pct" "e0021pct" "e0034pct" "e0055pct" "e0089pct")

for path in ${paths[@]}; do
    nohup go run github.com/jrhrmsll/stress@latest -throughput 15 -duration 3h -url http://localhost:8080/${path} &
done

exit 0
```

A hour later a second script should be executed to change the error rates as was described above. Finally, when all the alerts become triggered (around 30 minutes later) a last script execution changes the errors rates to the initial values.

All these scripts are include inside the directory `examples/errors-rates/scripts` using meaningful names:
- incident-start.sh
- incident-stop.sh
- run.sh

Both the "start" and "stop" scripts push a metric named `tsgen_incident_status` to Pushgateway. This is useful for some Grafana Dashboard to establish the beginning and end of errors.
