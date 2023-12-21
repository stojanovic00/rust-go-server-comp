#!/bin/bash

# starts servers and gets PID
export TCP_ADDRESS=localhost:8001
export THREAD_POOL_SIZE=8

sh -c 'echo $$ > get_pid.txt; exec  ../../../../rust_server/target/debug/rust_server ' &

# Waits for server to start
sleep 3

# Populates with data
curl --request PUT \
  --url http://localhost:8001/ \
  --header 'Content-Type: application/json' \
  --data '{ "id": 1, "description": "A", "value": 11 }'

#Track memory usage
pidstat  -r -p $(cat get_pid.txt) 1 > get_mem.txt &

#Track memory usage
pidstat  -u -p $(cat get_pid.txt) 1 > get_cpu.txt &

# Start benchmarking
ab  -c 100 -n 100000  http://localhost:8001/1 > get_bench.txt

kill $(cat get_pid.txt)


