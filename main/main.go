package main

import (
	"fmt"
	"net/http"
	"flag"
	"github.com/josephpanossian/urlshort"
)

func main() {
	var dataFormat = flag.String("dataF", "yaml", "The type of data input with format {url:key, value:path}, i.e. json, yaml}")

	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	if (*dataFormat == "yaml") {
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
		var err error
		handler, err = urlshort.YAMLHandler([]byte(yaml), handler)
		if err != nil {
			panic(err)
		}
	// Build JSONHandler with handler as fallback
	} else if (*dataFormat == "json") {
		json := `[
		{"path": "/josephpanossian", "url": "https://github.com/josephpanossian"},
		{"path": "/gh", "url": "https://github.com"}
		]
		`
		var err error
		handler, err  = urlshort.JSONHandler([]byte(json), handler)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
