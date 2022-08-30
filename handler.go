package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"github.com/boltdb/bolt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		if value, ok := pathsToUrls[r.URL.Path]; ok { 
			http.Redirect(w, r, value, http.StatusFound)	
		} else {
			fallback.ServeHTTP(w, r)	
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	var paths []pathYaml
	err := yaml.Unmarshal(yml, &paths)
	pathMap := make(map[string]string)
	for _, val := range paths {
		pathMap[val.Path] = val.URL
	}
	return MapHandler(pathMap, fallback), err
}

type pathYaml struct {
	Path string `yaml:"path"`
	URL string `yaml:"url"`
}

type pathJson struct {
	Path string
	URL string
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var paths []pathJson
	err := json.Unmarshal(jsn, &paths)
	pathMap := make(map[string]string)
	for _, val := range paths {
		pathMap[val.Path] = val.URL
	}
	return MapHandler(pathMap, fallback), err 
}

func DBHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	pathMap := map[string]string{}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("paths"))
		err := b.ForEach(func(k []byte, v []byte) error {
			pathMap[string(k)] = string(v)
			return nil
		})
		return err 
	})
	return MapHandler(pathMap, fallback), err
}