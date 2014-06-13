#!/bin/bash
bash run_experiment_aws.bash $1.10000.0 10000 0
sleep 10
bash run_experiment_aws.bash $1.100000.0 100000 0
sleep 10
bash run_experiment_aws.bash $1.1000000.0 1000000 0
sleep 10
bash run_experiment_aws.bash $1.10000000.0 10000000 0
sleep 10
bash run_experiment_aws.bash $1.100000000.0 100000000 0
sleep 10
bash run_experiment_aws.bash $1.10000.1000 10000 1000
sleep 10
bash run_experiment_aws.bash $1.100000.1000 100000 1000
sleep 10
bash run_experiment_aws.bash $1.1000000.1000 1000000 1000
sleep 10
bash run_experiment_aws.bash $1.10000000.1000 10000000 1000
sleep 10
bash run_experiment_aws.bash $1.100000000.1000 100000000 1000
