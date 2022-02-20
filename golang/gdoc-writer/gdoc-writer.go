package main

import (
	"gdoc-writer/google"
	"gdoc-writer/stdin"
	"google.golang.org/api/docs/v1"
	"log"
)

func main() {
	b, err := stdin.ReadStdin()
	if err != nil {
		log.Fatalf("Unable to read data from stdin: %v", err)
	}

	srv := google.GetDocumentService()

	newDoc := docs.Document{Title: "First document from Go API"}
	doc, err := srv.Documents.Create(&newDoc).Do()
	if err != nil {
		log.Fatalf("Unable to create document: %v", err)
	}

	md := docs.Dimension{Magnitude: 20.0, Unit: "PT"}
	style := docs.DocumentStyle{MarginLeft: &md, MarginRight: &md, MarginTop: &md}
	styleRequest := docs.UpdateDocumentStyleRequest{DocumentStyle: &style, Fields: "marginTop,marginLeft,marginRight"}

	ins1 := docs.InsertTextRequest{Text: string(b), Location: &docs.Location{Index: 1}}
	ins2 := docs.InsertTextRequest{Text: "Seba was here!\n", Location: &docs.Location{Index: 1}}

	requests := make([]*docs.Request, 0)
	requests = append(requests, &docs.Request{UpdateDocumentStyle: &styleRequest})
	requests = append(requests, &docs.Request{InsertText: &ins1})
	requests = append(requests, &docs.Request{InsertText: &ins2})
	_, err = srv.Documents.BatchUpdate(doc.DocumentId, &docs.BatchUpdateDocumentRequest{Requests: requests}).Do()
	if err != nil {
		log.Fatalf("Update to update document stype: %v", err)
	}

}
