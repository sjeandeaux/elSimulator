package main

import (
	"io"
	"net/http"
)

/**
* Request <= /folder0/..../folderN?param1=param1&...paramN=valueN with POST or PUT content
* Found file [folder configuration]/determine(request)
* If the file does not exist then it sends a 404 else the file is written to the stream
 */
func ElSimulator(
	w http.ResponseWriter,
	r *http.Request) {
	//TODO get file
	io.WriteString(w, "This life is a party!!!")
}

func main() {
	http.HandleFunc("/", ElSimulator)
	//TODO configuration binding
	http.ListenAndServe("localhost:4000", nil)
}
