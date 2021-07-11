for i in $(seq 1 10)
do 
	go run client_oneshot.go &
	sleep 2
done
