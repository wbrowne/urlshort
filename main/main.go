package main

import (
	"fmt"
	"github.com/wbrowne/urlshort"
	"io/ioutil"
	"net/http"
)

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlString := readFile("pathMapping.yml")
	yamlHandler, err := urlshort.YAMLHandler([]byte(yamlString), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonString := readFile("pathMapping.json")
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonString), yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func readFile(filename string) []byte {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return file
}
