package main

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/josephpanossian/urlshort"
)

func main() {
	// //input filename flag
	// var fName = flag.String("fName", "paths.yaml", "Specify file to parse with urls. If using json must also specify the -fType flag")
	// //could alternatively check extension type and use flag as fallback
	// var fType = flag.String("fType", "yaml", "Specify the type of file input. Supports yaml and json")
	// flag.Parse()
	db, err := bolt.Open("paths.db", 0600, nil)
	_ = db
	if err != nil {
		panic(err)
	}
	initDB(db)
	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := urlshort.MapHandler(pathsToUrls, mux)

	//read file into byte slice
	// data, err := os.ReadFile(*fName)
	// if err != nil {
	// 	panic(err)
	// }
	// if *fType == "yaml" { //YAML handler with default mux as fallback
	// 	fmt.Printf("serving with yaml file: %s\n", *fName)
	// 	handler, err = urlshort.YAMLHandler(data, handler)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else if *fType == "json" { //JSON handler with default mux as fallback
	// 	fmt.Printf("serving with json file: %s\n", *fName)
	// 	handler, err = urlshort.JSONHandler(data, handler)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	//start and serve
	handler, err = urlshort.DBHandler(db, handler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

// default connection
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// default connection func handler
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func initDB(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			return err
		}
		b.Put([]byte("/jp"), []byte("https://github.com/josephpanossian"))
		b.Put([]byte("/gh"), []byte("https://github.com"))
		b.Put([]byte("/google"), []byte("https://google.com"))
		return err
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
