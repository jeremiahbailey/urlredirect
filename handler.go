package urlredirect

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathURLYAML struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type pathURLJSON struct {
	Path string `json:"path"`
	URL  string `json:"URL"`
}

//buildMap creates a map[string]string from already parsed yaml or json 'blobs'
func buildMapfromYAML(pathUrls []pathURLYAML) map[string]string {
	pathsToURLs := make(map[string]string)
	for _, u := range pathUrls {
		pathsToURLs[u.Path] = u.URL
	}

	return pathsToURLs
}
func buildMapfromJSON(pathUrls []pathURLJSON) map[string]string {
	pathsToURLs := make(map[string]string)
	for _, u := range pathUrls {
		pathsToURLs[u.Path] = u.URL
	}

	return pathsToURLs
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToURLs map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToURLs[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)

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
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToURLs := buildMapfromYAML(parsedYaml)
	return MapHandler(pathsToURLs, fallback), nil
}

func parseYAML(data []byte) ([]pathURLYAML, error) {
	var pathUrls []pathURLYAML

	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return pathUrls, nil
}

//JSON Handler

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//  [
//    {
//		"path": "some-path",
//		"url": "https://example.com"
//	 },
//	 {
//		"path": "some-otherpath",
//		"url": "https://anothersite.com"
//	 }
// ]

func JSONHandler(jsonblob []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonblob)
	if err != nil {
		return nil, err
	}
	pathsToURLs := buildMapfromJSON(parsedJSON)

	return MapHandler(pathsToURLs, fallback), nil
}

func parseJSON(jsonblob []byte) ([]pathURLJSON, error) {
	var pathUrlsJSON []pathURLJSON
	err := json.Unmarshal(jsonblob, &pathUrlsJSON)
	if err != nil {
		return nil, fmt.Errorf("error unamarshal json %v", err)
	}

	return pathUrlsJSON, nil
}
