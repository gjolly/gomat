First, all the project must be built.
//Todo starting and using the API
Upon starting, a gossiper is created, along with a GUI, available on localhost/8080. From there, it is possible to monitor the gossiper, by changing the maximum computational power or adding peers. Also, all the tasks the node is computing are visible, and so is the list of peers.

## Run the GUI
To run the GUI, the deamon must be running. Then open in your internet browser the file `./Gossiper/GUI/homepage.html`.

## Run the unit tests
`go test -timeout 30s github.com\matei13\gomat\matrix -run ^TestAdd|TestSub|TestMul`
`go test -timeout 30s github.com\matei13\gomat\Daemon\gomatcore -run ^TestSplit|TestMerge|TestSplitMultAddMerge`