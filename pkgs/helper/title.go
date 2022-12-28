package helper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HtmlTitle(text string) string {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(text))

	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		title = strings.TrimSpace(doc.Find("meta[property='og:title']").AttrOr("content", ""))
	}
	if title == "" {
		title = strings.TrimSpace(doc.Find("meta[property='twitter:title']").AttrOr("content", ""))
	}
	title = strings.TrimSpace(title)
	return title
}
