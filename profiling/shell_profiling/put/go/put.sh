#!/bin/bash
# Arguments:
# 1: pool size
# 2: Number of requests
# 3: Number of connections

# Check if required arguments are provided
if [ $# -lt 3 ]; then
    echo "Usage: $0 <pool_size> <num_of_requests> <num_of_connections>"
    exit 1
fi


# starts servers and gets PID
export TCP_ADDRESS=localhost:8002
export THREAD_POOL_SIZE=$1

sh -c 'echo $$ > get_pid.txt; exec  ../../../../go_server/server' &

# Waits for server to start
sleep 3

#Track memory usage
pidstat  -r -p $(cat get_pid.txt) 1 > mem.txt &

#Track cpu usage
pidstat  -u -p $(cat get_pid.txt) 1 > cpu.txt &

# Start benchmarking
ab -n $2 -c $3 -T 'application/json' -u body.json http://localhost:8002/ > bench.txt

kill $(cat get_pid.txt)
rm get_pid.txt

# Parse all created data and analyze
./data_analyzer go $1 $2 $3

rm cpu.txt
rm mem.txt
rm bench.txt


