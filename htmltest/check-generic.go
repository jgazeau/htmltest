package htmltest

import (
	"fmt"

	"github.com/wjdp/htmltest/htmldoc"
	"github.com/wjdp/htmltest/issues"
	"golang.org/x/net/html"
)

// Checks the reference in the provided node and attribute key
func (hT *HTMLTest) checkGeneric(document *htmldoc.Document, node *html.Node, key string) {
	// Fail silently if attribute isn't present
	if !htmldoc.AttrPresent(node.Attr, key) {
		return
	}

	urlStr := htmldoc.GetAttr(node.Attr, key)
	ref := htmldoc.NewReference(document, node, urlStr)

	// Check attr isn't blank
	if urlStr == "" {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   fmt.Sprintf(node.Data, key, "is blank"),
			Reference: ref,
		})
	}

	// Check the reference
	hT.checkGenericRef(ref)
}

func (hT *HTMLTest) checkGenericRef(ref *htmldoc.Reference) {
	// Route reference check
	switch ref.Scheme() {
	case "http":
		hT.enforceHTTPS(ref)
		hT.checkExternal(ref)
	case "https":
		hT.checkExternal(ref)
	case "file":
		hT.checkInternal(ref)
	}
}

func (hT *HTMLTest) enforceHTTPS(ref *htmldoc.Reference) {
	// Does this url match an url ignore rule?
	if hT.opts.isURLIgnored(ref.URLString()) {
		return
	}

	// Does this url match an url include rule?
	if !hT.opts.isURLIncluded(ref.URLString()) {
		return
	}

	if hT.opts.EnforceHTTPS {
		hT.issueStore.AddIssue(issues.Issue{
			Level:     issues.LevelError,
			Message:   "is not an HTTPS target",
			Reference: ref,
		})
	}
}
