bash run_server_modprobe.bash $1.parrallel.1 3000 -numClients=1 "${@:2}"  
bash run_server_modprobe.bash $1.parrallel.2 3000 -numClients=2 "${@:2}"  
bash run_server_modprobe.bash $1.parrallel.5 3000 -numClients=5 "${@:2}"  
bash run_server_modprobe.bash $1.parrallel.5 3000 -numClients=10 "${@:2}"  
bash run_server_modprobe.bash $1.parrallel.5 3000 -numClients=20 "${@:2}"  
