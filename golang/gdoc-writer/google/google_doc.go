package google

import (
	"fmt"
	"google.golang.org/api/docs/v1"
	"log"
)

const DocumentLocationTemplate = "https://docs.google.com/document/d/%s/edit"

func CreateDocument(content string) string {
	srv := NewDocumentService()

	doc, err := srv.Documents.Create(&docs.Document{Title: "First document from Go API"}).Do()
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
	_, err = srv.Documents.BatchUpdate(doc.DocumentId, &docs.BatchUpdateDocumentRequest{Requests: requests}).Do()
	if err != nil {
		log.Fatalf("Update to update document stype: %v", err)
	}

	return fmt.Sprintf(DocumentLocationTemplate, doc.DocumentId)
}
