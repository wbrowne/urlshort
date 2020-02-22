package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
)

type UrlShort struct {
	Path string
	Url string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestPath := request.URL.Path

		if url, ok := pathsToUrls[requestPath]; ok {
			http.Redirect(writer, request, url, http.StatusSeeOther)
			return
		}

		fallback.ServeHTTP(writer, request)
	})
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYML(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedYaml, fallback), nil
}

func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(parsedJson, fallback), nil
}

func parseYML(yml []byte) (map[string]string, error) {
	var urlShort []UrlShort
	err := yaml.Unmarshal(yml, &urlShort)

	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	return createPathMappings(urlShort), nil
}

func parseJSON(yml []byte) (map[string]string, error) {
	var urlShort []UrlShort
	err := json.Unmarshal(yml, &urlShort)

	if err != nil {
		log.Fatalf("error: %v", err)
		return nil, err
	}

	return createPathMappings(urlShort), nil
}

func createPathMappings(urlShort []UrlShort) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, elem := range urlShort {
		pathsToUrls[elem.Path] = elem.Url
	}

	return pathsToUrls
}
