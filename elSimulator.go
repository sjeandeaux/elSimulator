package main

import (
	"bytes"
	"flag"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
)

//Configuration Application
type ElSimulatorConfig struct {
	//url to call
	bindingAddress string
	//directory with file to read
	baseDirectory string
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
		elSimulatorCurrent    = "elSimulatorCurrent"
		defaultBaseDirectory  = elSimulatorCurrent //elSimulatorCurrent if default is current directory
	)
	flag.StringVar(&elSimulatorConfig.bindingAddress, "bindingAddress", defaultBindingAddress, "The binding address")
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

	case elSimulatorCurrent:
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
//TODO administration status request in levelDb???
func main() {
	http.HandleFunc("/file/", ElSimulatorHandle)
	log.Println("start on ", elSimulatorConfig.bindingAddress)
	err := http.ListenAndServe(elSimulatorConfig.bindingAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Request <= /folder0/..../folderN?param1=param1&...paramN=valueN with POST or PUT content
// Found file [folder configuration]/determine(request)
// If the file does not exist then it sends a 404 else the file is written to the stream
//TODO add parameter in template
func ElSimulatorHandle(
	w http.ResponseWriter,
	r *http.Request) {
	f := findFile(r)
	if f == nil {
		http.Error(w, "The life is a party!!!", http.StatusNotFound)
	} else {
		//TODO template
		name := f.Name()
		t, err := template.ParseFiles(name)
		if err != nil {
			log.Printf("error template => %s", err)
		}
		w.Header().Add("content-type", mime.TypeByExtension(filepath.Ext(name)))
		t.Execute(w, nil)
	}

}

// Find file if not found (or a other error) nil else file.
// Url => localhost/file/test/sub?param1=value1&param2=value2
// File /.../file/test/sub/param1_value1_param2_value2
// If not found first file matchs pattern /.../file/test/sub/param1_value1_param2_value2*
func findFile(r *http.Request) *os.File {
	fileToRead := elSimulatorConfig.baseDirectory + strings.Replace(r.URL.Path, separatorURL, pathSeparator, -1) + pathSeparator + nameFile(r.URL.Query())
	log.Printf("file => %s", fileToRead)
	//file exists
	file, err := os.Open(fileToRead)
	if err == nil {
		return file
	}
	log.Printf("error => %s", err)
	allFile, errTwo := filepath.Glob(fileToRead + "*")
	if errTwo != nil || len(allFile) == 0 {
		log.Println(errTwo)
		return nil
	}

	log.Printf("file => %s", allFile[0])
	file, err = os.Open(allFile[0]) // For read access.
	if err != nil {
		return nil
	}
	return file

}

//query to generate name.
//TODO sort parameter and valueS
//TODO filter
//TODO body content.
func nameFile(query url.Values) string {
	var buffer bytes.Buffer
	//append all key value
	for key, value := range query {
		buffer.WriteString(separator)
		buffer.WriteString(key)
		buffer.WriteString(separator)
		buffer.WriteString(strings.Join(value, separator))
	}
	if buffer.Len() == 0 {
		return withoutParameter
	}
	//remove first separator
	return buffer.String()[1:]
}
