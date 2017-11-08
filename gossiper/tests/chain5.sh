init_path=$(pwd)
cd ../../build
go build ../gossiper
go build ../gossiper/client


./gossiper -UIPort=5004 -gossipPort=localhost:10004 -name=nodeE > $init_path/E.log &
./gossiper -UIPort=5003 -gossipPort=localhost:10003 -name=nodeD -peers=localhost:10004> $init_path/D.log &
./gossiper -UIPort=5002 -gossipPort=localhost:10002 -name=nodeC -peers=localhost:10003> $init_path/C.log &
./gossiper -UIPort=5001 -gossipPort=localhost:10001 -name=nodeB -peers=localhost:10002> $init_path/B.log &
./gossiper -UIPort=5000 -gossipPort=localhost:10000 -name=nodeA -peers=localhost:10001 > $init_path/A.log &

sleep 1
./client -UIPort=5000 -msg="1:A->D"
sleep 2
./client -UIPort=5004 -msg="2:D->A"
sleep 2

# clearing
killall gossiper
rm gossiper client
cd $init_path
