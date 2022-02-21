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
	if tag == nil {
		log.Fatalf("Unable to find psalm element id")
	}

	tagEm := html.GetTagByName(tag, "em")
	if tagEm == nil {
		log.Fatalf("Unable to find psalm chorus")
	}
	fmt.Println(html.RenderNodeContent(tagEm))

	content := html.FindAllParagraphsWithText(tag)
	for _, t := range content {
		fmt.Println()
		fmt.Println(html.StripHtmlTags(t))
	}
}
