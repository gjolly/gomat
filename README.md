# Gomat
## Feature
Decentralized matrix computation

## Build
```bash
$ go build
```
This build the Daemon

## Usage
 - Run the daemon and import the package gomat
 - Create new matrices with `gomat.New()`
 - Run any computation with `gomat.Add()`, `gomat.Mult()` and `gomat.subs()`

## Run the GUI
To run the GUI, the deamon must be running. Then open in your internet browser the file `./Gossiper/GUI/homepage.html`.

## Run the unit tests
`go test -timeout 30s github.com\matei13\gomat\matrix -run ^TestAdd|TestSub|TestMul`
`go test -timeout 30s github.com\matei13\gomat\Daemon\gomatcore -run ^TestSplit|TestMerge|TestSplitMultAddMerge`