#!/bin/bash
bash run_client.bash $1.parallel.1 $2 $3 --par 1
sleep 10
bash run_client.bash $1.parallel.2 $2 $3 --par 2
sleep 10
bash run_client.bash $1.parallel.5 $2 $3 --par 5
sleep 10
bash run_client.bash $1.parallel.10 $2 $3 --par 10
sleep 10
bash run_client.bash $1.parallel.20 $2 $3 --par 20
sleep 10
