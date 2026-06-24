// Package helpdoc owns human-facing CLI help text.
package helpdoc

// FlagDoc describes one command flag's help text.
type FlagDoc struct {
	UsageEN  string
	ReviewZH string
}

// CommandDoc describes public command help and its Chinese review text.
type CommandDoc struct {
	Key       string
	ShortEN   string
	LongEN    string
	ExampleEN string
	ReviewZH  string
	Flags     map[string]FlagDoc
}

// Docs returns every SeekMoon-authored command help document.
func Docs() map[string]CommandDoc {
	docs := englishDocs()
	for key, review := range chineseReviewDocs() {
		doc := docs[key]
		doc.ReviewZH = review.ReviewZH
		if doc.Flags == nil {
			doc.Flags = map[string]FlagDoc{}
		}
		for name, flag := range review.Flags {
			current := doc.Flags[name]
			current.ReviewZH = flag.ReviewZH
			doc.Flags[name] = current
		}
		docs[key] = doc
	}
	return docs
}

// ReviewDocsZH returns Chinese review text for authored command help.
func ReviewDocsZH() map[string]CommandDoc {
	return chineseReviewDocs()
}
