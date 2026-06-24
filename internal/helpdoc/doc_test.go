package helpdoc

import (
	"strings"
	"testing"
)

func TestDocsCoverPublicSeekMoonCommands(t *testing.T) {
	want := []string{
		"seekmoon",
		"seekmoon doctor",
		"seekmoon sync",
		"seekmoon search",
		"seekmoon view",
		"seekmoon api",
		"seekmoon source",
		"seekmoon skill",
		"seekmoon skill search",
		"seekmoon skill view",
		"seekmoon compare",
		"seekmoon probe",
		"seekmoon record",
		"seekmoon report",
		"seekmoon raw",
	}
	docs := Docs()
	reviews := ReviewDocsZH()
	for _, key := range want {
		doc, ok := docs[key]
		if !ok {
			t.Fatalf("missing doc %q", key)
		}
		if strings.TrimSpace(doc.ShortEN) == "" || strings.TrimSpace(doc.LongEN) == "" || strings.TrimSpace(doc.ExampleEN) == "" {
			t.Fatalf("incomplete English doc for %q: %#v", key, doc)
		}
		if strings.TrimSpace(doc.ReviewZH) == "" {
			t.Fatalf("missing merged Chinese review for %q", key)
		}
		if strings.TrimSpace(reviews[key].ReviewZH) == "" {
			t.Fatalf("missing Chinese review source for %q", key)
		}
	}
}

func TestCommonFlagDocsHaveEnglishAndChineseText(t *testing.T) {
	for _, name := range []string{"json", "jq", "shape", "schema"} {
		doc := CommonFlagDocs()[name]
		if strings.TrimSpace(doc.UsageEN) == "" || strings.TrimSpace(doc.ReviewZH) == "" {
			t.Fatalf("common flag %q doc = %#v", name, doc)
		}
	}
}
