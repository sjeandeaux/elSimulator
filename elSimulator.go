package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//Configuration Application
type ElSimulatorConfig struct {
	//url to call
	bindingAddress string
	//directory with file to read
	baseDirectory  string
	parameterRegex string
	proxyAddress   string
}

type Info struct {
	HttpCode       int `mmm`
	UrlRedirection string
	Header         map[string]string
}

const (
	prefixInfo       = "info_"
	suffixInfo       = ".json"
	separator        = "_"
	context          = "/file/"
	contextProxy     = "/proxy/"
	urlSeparator     = "/"
	withoutParameter = "withoutParameter"
	pathSeparator    = string(os.PathSeparator)
)

//configuration to use
var elSimulatorConfig = new(ElSimulatorConfig)

//parse command in configuration
func init() {
	const (
		defaultBindingAddress = "localhost:4000"
		defaultProxyAddress   = "http://localhost:4000/file"
		elSimulatorCurrent    = "elSimulatorCurrent"
		defaultBaseDirectory  = elSimulatorCurrent //elSimulatorCurrent if default is current directory
		defaultParameterRegex = ".*"
	)
	flag.StringVar(&elSimulatorConfig.bindingAddress, "bindingAddress", defaultBindingAddress, "The binding address")
	flag.StringVar(&elSimulatorConfig.proxyAddress, "proxyAddress", defaultProxyAddress, "The proxy address")
	flag.StringVar(&elSimulatorConfig.baseDirectory, "baseDirectory", defaultBaseDirectory, "directory with file to read (elSimulatorCurrent to use directory elSimulator)")
	flag.StringVar(&elSimulatorConfig.parameterRegex, "parameterRegex", defaultParameterRegex, "Parameter regex")
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

func main() {
	http.HandleFunc(contextProxy, ElProxyHandle)
	http.HandleFunc(context, ElSimulatorHandle)
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
	fileInfo, fileToRead, params := findFile(r)
	//TODO header in configuration
	addCORSHeader(w)
	log.Println(fileInfo)
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

		for key, value := range fileInfo.Header {
			w.Header().Add(key, value)
		}

		w.WriteHeader(fileInfo.HttpCode)
		t.Execute(w, params)

	}

}

//TODO error NewRequest and Do...
func ElProxyHandle(
	w http.ResponseWriter,
	r *http.Request) {
	addCORSHeader(w)
	base, calName, params := NameFile(r)
	indexPath := strings.Index(calName, pathSeparator)
	log.Println(base, calName, indexPath, params)
	calledUrl := elSimulatorConfig.proxyAddress + strings.Replace(r.URL.RequestURI(), contextProxy, urlSeparator, 1)
	log.Println("called url ", calledUrl)

	req, _ := http.NewRequest(r.Method, calledUrl, nil)
	resp, _ := http.DefaultClient.Do(req)

	if indexPath != -1 {
		os.MkdirAll(base+calName[:indexPath], 0755)
	} else {
		os.MkdirAll(base, 0755)
	}
	go saveInfo(resp, base, calName)
	file, errFile := os.Create(base + calName)
	defer file.Close()
	if errFile != nil {
		log.Println(base, calName, errFile)
	}

	for val := range resp.Header {
		log.Println("...", val, resp.Header.Get(val))
		w.Header().Set(val, resp.Header.Get(val))
	}
	w.WriteHeader(resp.StatusCode)
	// make a buffer to keep chunks that are read
	whatIRead := make([]byte, 1024)
	for {
		// read a chunk
		n, err := resp.Body.Read(whatIRead)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		_, errW := file.Write(whatIRead[:n])
		if errW != nil {
			log.Println(errW)
		}
		if _, err := w.Write(whatIRead[:n]); err != nil {

			panic(err)
		}
	}
}

//TODO configuration header
func addCORSHeader(w http.ResponseWriter) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS, PUT, DELETE, POST, PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
}

//save status and header in file info
//TODO error
func saveInfo(resp *http.Response, base, calName string) {
	var info Info
	info.HttpCode = resp.StatusCode
	info.Header = make(map[string]string)
	for val := range resp.Header {
		log.Println(val)
		info.Header[val] = resp.Header.Get(val)
	}
	infoJson, _ := json.MarshalIndent(info, "", "   ")

	file, errFile := os.Create(getFileNameInfo(base, calName))
	defer file.Close()
	if errFile != nil {
		log.Println(base, calName, errFile)
	}
	file.Write(infoJson)
	if errFile != nil {
		log.Println(base, calName, errFile)
	}

}

// Find file if not found (or a other error) nil else file.
// Url GET => localhost/file/test/sub?param1=value1&param2=value2
// File /.../file/test/sub/GET/param1_value1_param2_value2
// If not found first file matchs pattern /.../file/test/sub/GET/param1_value1_param2_value2*
func findFile(r *http.Request) (*Info, *os.File, map[string][]string) {

	base, calName, params := NameFile(r)

	fileToRead := base + calName
	fileInfo, errInfo := getInfo(base, calName)
	if errInfo != nil {
		log.Println(errInfo)
	}

	log.Printf("file => %s", fileToRead)
	//file exists
	file, err := os.Open(fileToRead)
	if err == nil {
		return fileInfo, file, params
	}
	log.Printf("error => %s", err)
	allFile, errTwo := filepath.Glob(fileToRead + "*")
	if errTwo != nil || len(allFile) == 0 {
		log.Println(errTwo)
		return fileInfo, nil, params
	}

	log.Printf("file => %s", allFile[0])
	file, err = os.Open(allFile[0]) // For read access.
	if err != nil {
		return fileInfo, nil, params
	}
	return fileInfo, file, params

}

func getFileNameInfo(base, calName string) string {
	return base + prefixInfo + calName + suffixInfo
}

//Read file json
func getInfo(base, calName string) (*Info, error) {
	fileToReadInfo := getFileNameInfo(base, calName)
	log.Println("file info ", fileToReadInfo)
	file, err := os.Open(fileToReadInfo)
	if err != nil {
		return &Info{404, "", nil}, err
	}
	bytesFile, errRead := ioutil.ReadAll(file)
	if errRead != nil {
		log.Println("error:", errRead)
		return &Info{200, "", nil}, nil
	}
	var info Info
	errJson := json.Unmarshal(bytesFile, &info)
	if errJson != nil {
		log.Println("error:", errJson)
		return &Info{200, "", nil}, nil
	}
	return &info, nil
}

type NameFileParameter struct {
	buffer       bytes.Buffer
	allParameter map[string][]string
}

func (b *NameFileParameter) GetName() string {
	if b.buffer.Len() == 0 {
		return withoutParameter
	}
	//remove first separator
	return b.buffer.String()[1:]
}

// Query parses RawQuery and returns the corresponding values.
func (b *NameFileParameter) Append(values url.Values) {
	for key, value := range values {
		b.allParameter[key] = value
		if match, _ := regexp.Match(elSimulatorConfig.parameterRegex, []byte(key)); match {
			log.Println("key ", key)
			b.buffer.WriteString(separator)
			b.buffer.WriteString(key)
			b.buffer.WriteString(separator)
			b.buffer.WriteString(strings.Join(value, separator))
		}

	}
}

func (b *NameFileParameter) Base(method, path string) string {
	var folder string
	if path == context {
		folder = strings.Replace(path[:len(context)-1], urlSeparator, pathSeparator, -1)

	} else {
		folder = strings.Replace(path, urlSeparator, pathSeparator, -1)
	}

	return elSimulatorConfig.baseDirectory + folder + pathSeparator + method + pathSeparator
}

//query to generate name.
func NameFile(r *http.Request) (string, string, map[string][]string) {
	var buffer NameFileParameter
	buffer.allParameter = make(map[string][]string)
	//append all key value
	buffer.Append(r.URL.Query())
	buffer.Append(r.Form)
	return buffer.Base(r.Method, r.URL.Path), buffer.GetName(), buffer.allParameter
}
