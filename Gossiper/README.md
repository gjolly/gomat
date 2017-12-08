# Gossiper
A gossip network like Usenet.

## How to
### Compile
Just run: `go build` in the main directory.

>You will need those two libraries:
> - `github.com/dedis/protobuf`
> - `github.com/gorilla/mux`
>
>Install them with `go get LIBRARY` 

 
### Run a node
Use Gossiper:
```
Usage of Gossiper:
  -UIPort string
    	UIPort (default "10000")
  -gossipAddr string
    	gossipAddr:gossipPort (default "localhost:5000")
  -guiAddr string
    	Address of the GUI: guiAddr:guiPort (default "none")
  -name string
    	Name of the node (default "nodeA")
  -peers string
    	List of peers: addrPeer1:portPeer1_addrPeer2:portPeer2 ...
  -rtimer uint
    	Delay during two routing message (Developer) (default 60)
```

> If you want to be part of a network, you have to specify at least one peer


### Use the GUI
Open a web browser and go on guiAddr:guiPort.
