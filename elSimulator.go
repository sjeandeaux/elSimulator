package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

/**
* Configuration Application
 */
type ElSimulatorConfig struct {
	//url to call
	bindingAddress string
}

//configuration to use
var elSimulatorConfig = new(ElSimulatorConfig)

//parse command in configuration
func init() {
	flag.StringVar(&elSimulatorConfig.bindingAddress, "bindingAddress", "localhost:4000", "The binding address")
	flag.Parse()
}

/**
* Bind address.
 */
func main() {
	http.HandleFunc("/", ElSimulatorHandle)
	log.Println("start on %s", elSimulatorConfig.bindingAddress)
	err := http.ListenAndServe(elSimulatorConfig.bindingAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

/**
* Request <= /folder0/..../folderN?param1=param1&...paramN=valueN with POST or PUT content
* Found file [folder configuration]/determine(request)
* If the file does not exist then it sends a 404 else the file is written to the stream
 */
func ElSimulatorHandle(
	w http.ResponseWriter,
	r *http.Request) {
	//TODO get file
	io.WriteString(w, "This life is a party!!!")

}
