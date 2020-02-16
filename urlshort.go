package urlshort

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

type UrlShort struct {
	Path string
	Url string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestPath := request.URL.Path

		url := pathsToUrls[requestPath]

		if url != "" {
			http.Redirect(writer, request, url, http.StatusSeeOther)
			return
		}

		fallback.ServeHTTP(writer, request)
	})
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
	parsedYaml, err := parseYML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedYaml, fallback), nil
}

func parseYML(yml []byte) (map[string]string, error) {
	var urlShort []UrlShort
	err := yaml.Unmarshal(yml, &urlShort)

	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	pathsToUrls := make(map[string]string)

	for _, elem := range urlShort {
		pathsToUrls[elem.Path] = elem.Url
	}

	return pathsToUrls, nil
}
