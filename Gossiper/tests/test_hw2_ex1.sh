#!/usr/bin/env bash
set -e

cd ../

go build
cd client
go build
cd ..
mv Gossiper gossiper

RED='\033[0;31m'
NC='\033[0m'
DEBUG="false"

outputFiles=()
message=Weather_is_clear
message2=Winter_is_coming

# Making a simple network:
#      A
#     / \
#    B  C
#   / \
#  D  E
startGossip(){
	local name=$1
	local port=$2
	local peers=""
	if [ "$3" ]; then
		peers="-peers=127.0.0.1:$3"
	fi
	echo ./gossiper -gossipAddr=127.0.0.1:$port -UIPort=$((port+1)) -name=$name $peers
	./gossiper -gossipAddr=127.0.0.1:$port -UIPort=$((port+1)) -name=$name $peers> $name.log &
	# don't show 'killed by signal'-messages
	disown
}
startGossip A 10000
startGossip B 10002 10000
startGossip C 10004 10000
startGossip D 10006 10002
startGossip E 10008 10002
sleep 1

# Give message to E
./client/client -UIPort=10009 -msg=$message
sleep 1

# Give message to D
./client/client -UIPort=10007 -msg=$message2
# Give message to A
./client/client -UIPort=10001 -msg=$message2
sleep 2
pkill -f gossiper

#testing
fail(){
	echo -e "${RED}*** Failed test $1 ***${NC}"
  exit 1
}

grep -q "DSDV E:127.0.0.1:10002" A.log || fail "B doesn't see D through B"
grep -q "DSDV D:127.0.0.1:10006" B.log || fail "B doesn't see D through B"
grep -q "DSDV E:127.0.0.1:10008" B.log || fail "B doesn't see E through B"

rm gossiper client/client
