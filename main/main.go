package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jeremiahbailey/urlredirect"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlredirect.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile("../inputs/urls.yaml")
	if err != nil {
		fmt.Printf("unable to parse yaml due to %v", err)
	}

	yamlHandler, err := urlredirect.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the mapHandler as the
	// fallback
	json, err := ioutil.ReadFile("../inputs/urls.json")
	if err != nil {
		fmt.Printf("error reading json file: %v", err)
	}

	jsonHandler, err := urlredirect.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8000")
	http.ListenAndServe(":8000", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
