#!/usr/bin/env bash

HOST="localhost:8080"

curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0002pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0003pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0005pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0008pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0013pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0021pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0034pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0055pct/faults/500 --data '{"rate": 0.001}'
curl -X POST -H "Content-Type: application/json" http://${HOST}/api/paths/e0089pct/faults/500 --data '{"rate": 0.001}'

cat <<EOF | curl --data-binary @- http://localhost:9091/metrics/job/tsgen/instance/bash
# TYPE tsgen_incident_status gauge
tsgen_incident_status{} 0
EOF

exit 0
