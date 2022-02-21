package main

import (
	"fmt"
	"log"
	"net/http"
	"psalm-fetcher/html"
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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Unable to parse HTML document: %v", err)
	}

	tag := html.GetElementById(doc, "tabnowy01")
	if err != nil {
		log.Fatalf("Unable to find element with id tabnowy01: %v", err)
	}

	fmt.Println(html.RenderNode(tag))

}
