/*
Package html helps search through HTML document and find elements by ID, tags by name, display nodes, etc.
*/
package html

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"strings"
	"unicode/utf8"
)

// Parse returns html node that could be used in functions from this package or error if issues when parsing
func Parse(reader io.Reader) (*html.Node, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// GetElementById search for id in given html fragment and returns html node when element was found
// It returns nil otherwise.
func GetElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id, "")
}

// GetTagByName search for given tag name in given html fragment and returns html node when element was found.
// It returns nil otherwise.
func GetTagByName(n *html.Node, name string) *html.Node {
	return traverse(n, "", name)
}

func FindAllParagraphsWithText(n *html.Node) []string {
	output := make([]string, 0)
	if content := paragraphWithContent(n); content != "" {
		return append(output, content)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content := FindAllParagraphsWithText(c)
		if len(content) > 0 {
			output = append(output, content...)
		}
	}

	return output
}

// RenderNode renders given html fragment and returns string
func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)
	if err != nil {
		return ""
	}

	return buf.String()
}

func RenderNodeContent(n *html.Node) string {
	return n.FirstChild.Data
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

func checkName(n *html.Node, name string) bool {
	return n.Type == html.ElementNode && n.Data == name
}

func traverse(n *html.Node, id string, name string) *html.Node {
	if id != "" && name != "" && checkId(n, id) && checkName(n, name) {
		return n
	} else if id != "" && checkId(n, id) {
		return n
	} else if name != "" && checkName(n, name) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := traverse(c, id, name)
		if res != nil {
			return res
		}
	}

	return nil
}

func paragraphWithContent(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "p" && n.FirstChild.Type == html.TextNode {
		return RenderNode(n)
	}
	return ""
}

const (
	htmlTagStart = 60 // Unicode `<`
	htmlTagEnd   = 62 // Unicode `>`
)

// Aggressively strips HTML tags from a string.
// It will only keep anything between `>` and `<`.
func StripHtmlTags(s string) string {
	// Setup a string builder and allocate enough memory for the new string.
	var builder strings.Builder
	builder.Grow(len(s) + utf8.UTFMax)

	in := false // True if we are inside an HTML tag.
	start := 0  // The index of the previous start tag character `<`
	end := 0    // The index of the previous end tag character `>`

	for i, c := range s {
		// If this is the last character and we are not in an HTML tag, save it.
		if (i+1) == len(s) && end >= start {
			builder.WriteString(s[end:])
		}

		// Keep going if the character is not `<` or `>`
		if c != htmlTagStart && c != htmlTagEnd {
			continue
		}

		if c == htmlTagStart {
			// Only update the start if we are not in a tag.
			// This make sure we strip out `<<br>` not just `<br>`
			if !in {
				start = i
			}
			in = true

			// Write the valid string between the close and start of the two tags.
			builder.WriteString(s[end:start])
			continue
		}
		// else c == htmlTagEnd
		in = false
		end = i + 1
	}
	s = builder.String()
	return s
}
