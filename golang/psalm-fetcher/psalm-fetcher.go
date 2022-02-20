package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// psalmSourceUrlTemplate is a template URL from where we fetch psalm content
	psalmSourceUrlTemplate = "https://niezbednik.niedziela.pl/liturgia/%s/Psalm"
)

func main() {
	currentTime := time.Now()
	resp, err := http.Get(fmt.Sprintf(psalmSourceUrlTemplate, currentTime.Format("2006-01-02")))
	if err != nil {
		log.Fatalf("Unable to GET content from site: %v", err)
	}
	defer resp.Body.Close()

	//c, err := io.ReadAll(resp.Body)
	//println(string(c))

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Unable to parse document: %v", err)
	}

	tag := getElementById(doc, "tabnowy01")
	fmt.Println(renderNode(tag))
}

func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id)
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}

	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "id")

		if ok && s == id {
			return true
		}
	}

	return false
}

func traverse(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := traverse(c, id)
		if res != nil {
			return res
		}
	}

	return nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)
	if err != nil {
		return ""
	}

	return buf.String()
}
