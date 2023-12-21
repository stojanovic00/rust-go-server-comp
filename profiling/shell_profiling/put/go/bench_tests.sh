#!/bin/bash
# PARAMS: pool_size requests connections
./put.sh 2 100000 10
./put.sh 4 100000 10
./put.sh 8 100000 10
./put.sh 16 100000 10
./put.sh 32 100000 10
./put.sh 64 100000 10

./put.sh 2 100000 100
./put.sh 4 100000 100
./put.sh 8 100000 100
./put.sh 16 100000 100
./put.sh 32 100000 100
./put.sh 64 100000 100

./put.sh 2 100000 1000
./put.sh 4 100000 1000
./put.sh 8 100000 1000
./put.sh 16 100000 1000
./put.sh 32 100000 1000
./put.sh 64 100000 1000
