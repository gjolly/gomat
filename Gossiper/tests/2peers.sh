init_path=$(pwd)
cd ../build
go build ../
go build ../client


./Gossiper -UIPort=5001 -gossipPort=localhost:10001 -name=nodeB > $init_path/B.log &
./Gossiper -UIPort=5000 -gossipPort=localhost:10000 -name=nodeA -peers=localhost:10001 > $init_path/A.log &

sleep 1
./client -UIPort=5000 -msg="A->B"
sleep 1

# clearing
killall Gossiper
rm Gossiper client
cd $init_path

# Analyse
GREEN='\033[0;32m'
if grep -q "A->B" B.log ; then
	echo -e "${GREEN}A->B succed"
fi
