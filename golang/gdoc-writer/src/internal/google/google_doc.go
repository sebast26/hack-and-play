package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/docs/v1"
)

const (
	// documentLocationTemplate is a template used by Google Docs to access document by ID from the browser
	documentLocationTemplate = "https://docs.google.com/document/d/%s/edit"

	// maxContentLengthToTitle specifies maximum number of characters that could be included inside document title
	maxContentLengthToTitle = 30
)

// CreateDocument creates Google Document with context and prefix
func CreateDocument(service *docs.Service, content string, prefix string) string {
	title := createDocumentTitle(content, prefix)
	doc, err := service.Documents.Create(&docs.Document{Title: title}).Do()
	if err != nil {
		log.Fatalf("Unable to create document: %v", err)
	}

	md := docs.Dimension{Magnitude: 20.0, Unit: "PT"}
	style := docs.DocumentStyle{MarginLeft: &md, MarginRight: &md, MarginTop: &md}
	styleRequest := docs.UpdateDocumentStyleRequest{DocumentStyle: &style, Fields: "marginTop,marginLeft,marginRight"}

	ins1 := docs.InsertTextRequest{Text: content, Location: &docs.Location{Index: 1}}

	requests := make([]*docs.Request, 0)
	requests = append(requests, &docs.Request{UpdateDocumentStyle: &styleRequest})
	requests = append(requests, &docs.Request{InsertText: &ins1})
	_, err = service.Documents.BatchUpdate(doc.DocumentId, &docs.BatchUpdateDocumentRequest{Requests: requests}).Do()
	if err != nil {
		log.Fatalf("Update to update document stype: %v", err)
	}

	return fmt.Sprintf(documentLocationTemplate, doc.DocumentId)
}

func createDocumentTitle(context string, prefix string) string {
	i := strings.Index(context, "\n")
	if i == -1 || i > maxContentLengthToTitle {
		i = maxContentLengthToTitle
	}

	currentTime := time.Now()
	if prefix == "" {
		return fmt.Sprintf("%s - %s...", currentTime.Format("2006-01-02"), context[:i])
	}
	return fmt.Sprintf("%s %s - %s...", prefix, currentTime.Format("2006-01-02"), context[:i])
}
