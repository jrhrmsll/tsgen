#!/usr/bin/env bash

paths=("e0002pct" "e0003pct" "e0005pct" "e0008pct" "e0013pct" "e0021pct" "e0034pct" "e0055pct" "e0089pct")

for path in ${paths[@]}; do
    nohup go run github.com/jrhrmsll/stress@latest -throughput 15 -duration 4h -url http://localhost:8080/${path} &
done

exit 0
