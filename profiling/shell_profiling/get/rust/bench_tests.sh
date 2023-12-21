#!/bin/bash
# PARAMS: pool_size requests connections
./get.sh 2 100000 10
./get.sh 4 100000 10
./get.sh 8 100000 10
./get.sh 16 100000 10
./get.sh 32 100000 10
./get.sh 64 100000 10

./get.sh 2 100000 100
./get.sh 4 100000 100
./get.sh 8 100000 100
./get.sh 16 100000 100
./get.sh 32 100000 100
./get.sh 64 100000 100

./get.sh 2 100000 1000
./get.sh 4 100000 1000
./get.sh 8 100000 1000
./get.sh 16 100000 1000
./get.sh 32 100000 1000
./get.sh 64 100000 1000
