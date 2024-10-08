package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"parser/parser"
)

func main() {
	stack := make([]string, 0, 10)
	m := make(map[string]parser.Link)
	stack = append(stack, "https://pkg.go.dev/net/http")
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res, err := http.Get(u)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return
		}
		body, err := io.ReadAll(res.Body)
		r := bytes.NewReader(body)
		links := parser.GetLinks(r)
		for _, link := range links {
			u, e := url.Parse(link.Href)
			if e != nil {
				continue
			}
			if u.Host == "" {
				link.Href = "https://pkg.go.dev" + u.Path
			}
			if _, ok := m[link.Href]; !ok {
				stack = append(stack, link.Href)
				m[link.Href] = link
			}
		}
	}
}
