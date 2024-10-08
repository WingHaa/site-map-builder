package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"parser/parser"
	"path"
	"strings"
)

func formatDomain(domain string) string {
	return strings.TrimPrefix(domain, "www.")
}

func main() {
	p := "https://calhoun.io"
	u, e := url.Parse(p)
	if e != nil {
		panic(e)
	}
	origin := formatDomain(u.Hostname())
	fmt.Println("domain:", origin)
	stack := make([]string, 0, 10)
	m := make(map[string]parser.Link)
	stack = append(stack, "https://calhoun.io")
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		res, err := http.Get(u)
		if err != nil {
			fmt.Println("Error fetching URL", err)
			continue
		}
		if res.StatusCode != http.StatusOK {
			continue
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		r := bytes.NewReader(body)
		links := parser.GetLinks(r)
		for _, link := range links {
			u, e := url.Parse(link.Href)
			if e != nil || u.Path == "" {
				continue
			}

			domain := formatDomain(u.Hostname())
			if domain != "" && domain != origin {
				continue
			}

			p := link.Href
			if domain == "" {
				p = strings.TrimRight("https://"+path.Join(origin, u.Path), "/")
			}

			if _, ok := m[p]; !ok {
				fmt.Println("Wrote", p)
				stack = append(stack, p)
				m[p] = link
			}
		}
	}

	fmt.Print(m)
}
