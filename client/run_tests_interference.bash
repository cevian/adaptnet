#!/bin/bash
mkdir data/$1
go run client.go -addr sns58.cs.princeton.edu:3001 -bitsPerChunk=100000000 -numChunks 100 -msBetweenChunks 0 2>&1| tee data/$1/client.out
##bash run_experiment_aws.bash $1_interference_1 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_2 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_3 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_4 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_5 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_6 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_7 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_8 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_9 100000000 0
#sleep 10
#bash run_experiment_aws.bash $1_interference_10 100000000 0
#sleep 10
