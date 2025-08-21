package buckets

import (
	"net/url"
	"strconv"
)

// Location 表示解析后的 URL
type Location struct {
	Scheme   string
	User     string
	Host     string
	Port     int
	Path     string
	Query    map[string]string
	Fragment string
}

// location is a url string
func ParseLocation(location string) (*Location, error) {

	u, err := url.ParseRequestURI(location)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	portStr := u.Port()
	portInt, _ := strconv.Atoi(portStr)

	l := new(Location)
	l.Scheme = u.Scheme
	l.User = u.User.Username()
	l.Host = u.Hostname()
	l.Port = portInt
	l.Path = u.Path
	l.Fragment = u.Fragment
	l.Query = make(map[string]string)

	for key, vlist := range query {
		for _, value := range vlist {
			l.Query[key] = value
		}
	}

	return l, nil
}
