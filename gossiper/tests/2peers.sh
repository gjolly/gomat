init_path=$(pwd)
cd ../../build
go build ../gossiper
go build ../gossiper/client


./gossiper -UIPort=5001 -gossipPort=localhost:10001 -name=nodeB > $init_path/B.log &
./gossiper -UIPort=5000 -gossipPort=localhost:10000 -name=nodeA -peers=localhost:10001 > $init_path/A.log &

sleep 1
./client -UIPort=5000 -msg="1:A->B"
sleep 1

# clearing
killall gossiper
rm gossiper client
cd $init_path
