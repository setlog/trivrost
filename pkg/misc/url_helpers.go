package misc

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

func MustJoinURL(base, reference string) string {
	baseURL := MustParseURL(base)
	baseURL.Path = path.Join(baseURL.Path, filepath.ToSlash(reference))
	return baseURL.String()
}

func MustStripLastURLPathElement(urlString string) string {
	baseURL := MustParseURL(urlString)
	baseURL.Path = path.Dir(baseURL.Path)
	return baseURL.String()
}

func MustParseURL(urlString string) *url.URL {
	baseURL, baseErr := url.Parse(urlString)
	if baseErr != nil {
		panic(fmt.Sprintf("Could not parse URL \"%s\": %v", urlString, baseErr))
	}
	return baseURL
}

func EllipsisURL(urlString string, truncateThreshold int) string {
	baseURL, baseErr := url.Parse(urlString)
	if baseErr != nil {
		return urlString
	}
	if len(urlString) >= truncateThreshold {
		base := path.Base(baseURL.Path)
		if strings.HasSuffix(baseURL.Path, "/") && strings.Count(baseURL.Path, "/") < len(baseURL.Path) {
			base += "/"
		}
		if base == "/" {
			base = ""
		}
		ellipsisURL := baseURL.Scheme + "://" + baseURL.Host + "/.../" + base
		if len(ellipsisURL) < len(urlString) {
			return ellipsisURL
		}
	}
	return urlString
}
