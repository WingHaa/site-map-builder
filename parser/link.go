package parser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func isLinkTag(z *html.Tokenizer) bool {
	tn, _ := z.TagName()
	return len(tn) == 1 && tn[0] == 'a'
}

func getTagVal(t *html.Tokenizer, tag string) string {
	for {
		attr, val, more := t.TagAttr()
		if string(attr) == tag {
			return string(val)
		}
		if !more {
			return ""
		}
	}
}

func GetLinks(r io.Reader) []Link {
	links := make([]Link, 0, 5)
	tokenizer := html.NewTokenizer(r)
	var text strings.Builder
	var link Link
	depth := 0

HtmlLoop:
	for {
		token := tokenizer.Next()
		if err := tokenizer.Err(); err != nil {
			break
		}

		switch token {
		case html.CommentToken:
		case html.DoctypeToken:
		case html.SelfClosingTagToken:
			continue
		case html.ErrorToken:
			break HtmlLoop
		case html.TextToken:
			if depth == 1 {
				text.WriteString(strings.TrimSpace(string(tokenizer.Text())))
			}
		case html.StartTagToken:
			if !isLinkTag(tokenizer) {
				continue
			}

			if depth == 0 {
				link.Href = getTagVal(tokenizer, "href")
			}
			depth++
		case html.EndTagToken:
			if !isLinkTag(tokenizer) {
				continue
			}

			if depth == 1 {
				link.Text = text.String()
				links = append(links, link)
				text.Reset()
				link = Link{}
			}
			depth--
		default:
			panic(fmt.Sprintf("unexpected html.TokenType: %#v", token))
		}
	}

	return links
}
