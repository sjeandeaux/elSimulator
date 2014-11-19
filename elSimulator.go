package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
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

type Info struct {
	HttpCode       int
	UrlRedirection string
}

const (
	prefixInfo       = "info_"
	suffixInfo       = ".json"
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
//TODO Protocole CORS
func ElSimulatorHandle(
	w http.ResponseWriter,
	r *http.Request) {
	fileInfo, fileToRead := findFile(r)
	if fileToRead == nil {
		http.Error(w, "The life is a party!!!", fileInfo.HttpCode)
	} else {
		//TODO template
		name := fileToRead.Name()
		t, err := template.ParseFiles(name)
		if err != nil {
			log.Printf("error template => %s", err)
		}
		//read information status code and other
		w.Header().Add("content-type", mime.TypeByExtension(filepath.Ext(name)))
		w.WriteHeader(fileInfo.HttpCode)
		t.Execute(w, nil)

	}

}

// Find file if not found (or a other error) nil else file.
// Url => localhost/file/test/sub?param1=value1&param2=value2
// File /.../file/test/sub/param1_value1_param2_value2
// If not found first file matchs pattern /.../file/test/sub/param1_value1_param2_value2*
func findFile(r *http.Request) (*Info, *os.File) {
	base := elSimulatorConfig.baseDirectory + strings.Replace(r.URL.Path, separatorURL, pathSeparator, -1) + pathSeparator
	calName := nameFile(r.URL.Query())
	fileToRead := base + calName
	fileInfo, errInfo := getInfo(base, calName)
	if errInfo != nil {
		log.Println(errInfo)
	}

	log.Printf("file => %s", fileToRead)
	//file exists
	file, err := os.Open(fileToRead)
	if err == nil {
		return fileInfo, file
	}
	log.Printf("error => %s", err)
	allFile, errTwo := filepath.Glob(fileToRead + "*")
	if errTwo != nil || len(allFile) == 0 {
		log.Println(errTwo)
		return fileInfo, nil
	}

	log.Printf("file => %s", allFile[0])
	file, err = os.Open(allFile[0]) // For read access.
	if err != nil {
		return fileInfo, nil
	}
	return fileInfo, file

}

//Read file json
func getInfo(base, calName string) (*Info, error) {
	fileToReadInfo := base + prefixInfo + calName + suffixInfo
	log.Println("file info ", fileToReadInfo)
	file, err := os.Open(fileToReadInfo)
	if err != nil {
		return &Info{404, ""}, err
	}
	bytesFile, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		log.Println("error:", errRead)
		return &Info{200, ""}, nil
	}

	var info Info
	errJson := json.Unmarshal(bytesFile, &info)
	if errJson != nil {
		log.Println("error:", errJson)
		return &Info{200, ""}, nil
	}
	return &info, nil
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
