package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

//Configuration Application
type ElSimulatorConfig struct {
	//url to call
	bindingAddress string
	//directory with file to read
	baseDirectory string
	//context web
	context string
}

const (
	separator        = "_"
	separatorURL     = "/"
	withoutParameter = "withoutParameter"
	pathSeparator    = string(os.PathSeparator)
)

//configuration to use
var elSimulatorConfig = new(ElSimulatorConfig)

//parse command in configuration
func init() {
	const (
		defaultBindingAddress = "localhost:4000"
		defaultContext        = "/elSimulator/"
		defaultBaseDirectory  = "" //elSimulatorCurrent if default is current directory
	)
	flag.StringVar(&elSimulatorConfig.bindingAddress, "bindingAddress", defaultBindingAddress, "The binding address")
	flag.StringVar(&elSimulatorConfig.context, "context", defaultContext, "The context")
	flag.StringVar(&elSimulatorConfig.baseDirectory, "baseDirectory", defaultBaseDirectory, "directory with file to read (elSimulatorCurrent to use directory elSimulator)")
	flag.Parse()
	//use home'user
	switch elSimulatorConfig.baseDirectory {
	case "":
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		elSimulatorConfig.baseDirectory = usr.HomeDir
		break

	case "elSimulatorCurrent":
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		elSimulatorConfig.baseDirectory = dir
		break
	}
	log.Println("configuration :", elSimulatorConfig)
}

//Bind address.
func main() {
	http.HandleFunc(elSimulatorConfig.context, ElSimulatorHandle)
	log.Println("start on ", elSimulatorConfig.bindingAddress)
	err := http.ListenAndServe(elSimulatorConfig.bindingAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Request <= /folder0/..../folderN?param1=param1&...paramN=valueN with POST or PUT content
// Found file [folder configuration]/determine(request)
// If the file does not exist then it sends a 404 else the file is written to the stream
func ElSimulatorHandle(
	w http.ResponseWriter,
	r *http.Request) {
	f := findFile(r)
	if f == nil {
		http.Error(w, "The life is a party!!!", http.StatusNotFound)
	} else {
		//TODO template
		io.WriteString(w, "With file, it should be better")
	}

}

// Find file if not found (or a other error) nil else file.
func findFile(r *http.Request) *os.File {
	log.Println(elSimulatorConfig.baseDirectory, r.URL.Path, r.URL.Query(), r.Form)
	//TODO for on type file .txt, .xml, .json
	fileToRead := elSimulatorConfig.baseDirectory + strings.Replace(r.URL.Path, separatorURL, pathSeparator, -1) + pathSeparator + nameFile(r.URL.Query())
	log.Println(fileToRead)
	file, err := os.Open(fileToRead) // For read access.
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

//query to generate name.
//TODO sort
//TODO filter
//TODO body content.
func nameFile(query url.Values) string {
	var buffer bytes.Buffer

	for key, value := range query {
		buffer.WriteString(separator)
		buffer.WriteString(key)
		buffer.WriteString(separator)
		buffer.WriteString(strings.Join(value, separator))
	}
	if buffer.Len() == 0 {
		return withoutParameter
	}
	return buffer.String()[1:]
}
