package source

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MichengLiang/seekmoon/internal/model"
)

func TestAssetModuleIndexChildsAndRelpath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"root","package":{"path":"moonbitlang/core"},"childs":[{"name":"argparse","package":{"path":"moonbitlang/core/argparse","values":[{"name":"parse"}]}}]}`))
	}))
	defer server.Close()

	client := AssetClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchModuleIndex(context.Background(), "moonbitlang/core", "0.1.0")
	if result.Status != model.StatePresent {
		t.Fatalf("module index result = %#v", result)
	}
	if result.Value == nil || len(result.Value.Childs) != 1 {
		t.Fatalf("tree = %#v", result.Value)
	}
	rootRel := result.Value.Package.RelPath.Value
	childRel := result.Value.Childs[0].Package.RelPath.Value
	if rootRel == nil || *rootRel != "" || childRel == nil || *childRel != "argparse" {
		t.Fatalf("relpaths root=%v child=%v", rootRel, childRel)
	}
}

func TestAssetPackageDataPreservesRawAndPlainSignature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/argparse/package_data.json") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"argparse","values":[{"name":"parse","docstring":"","signature":"<a>parse</a>(x : String) -&gt; Unit","loc":{"file":"a.mbt","line":1}}]}`))
	}))
	defer server.Close()

	client := AssetClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchPackageData(context.Background(), "moonbitlang/core", "0.1.0", "moonbitlang/core/argparse")
	if result.Status != model.StatePresent || result.Value == nil || len(result.Value.Values) != 1 {
		t.Fatalf("package data result = %#v", result)
	}
	entry := result.Value.Values[0]
	if entry.Signature != "<a>parse</a>(x : String) -&gt; Unit" {
		t.Fatalf("raw signature changed: %q", entry.Signature)
	}
	if entry.PlainSignature.Value == nil || *entry.PlainSignature.Value != "parse(x : String) -> Unit" {
		t.Fatalf("plain signature = %#v", entry.PlainSignature)
	}
	if entry.Docstring.Status != model.StateMissing || entry.Loc.Status != model.StatePresent {
		t.Fatalf("entry evidence = %#v", entry)
	}
}

func TestAssetPackageDataNestedRelpathPreservesPathSegments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/pkg/sub/package_data.json") {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"sub"}`))
	}))
	defer server.Close()

	client := AssetClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchPackageData(context.Background(), "moonbitlang/core", "0.1.0", "moonbitlang/core/pkg/sub")
	if result.Status != model.StatePresent {
		t.Fatalf("package data result = %#v", result)
	}
}

func TestAssetResource404Unavailable(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	client := AssetClient{BaseURL: server.URL, Fetcher: Fetcher{Client: server.Client()}}
	result := client.FetchResource(context.Background(), "moonbitlang/core", "0.1.0", "moonbitlang/core")
	if result.Status != model.StateUnavailable {
		t.Fatalf("resource 404 status = %s", result.Status)
	}
}

func TestSourceZipSummary(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, name := range []string{"moon.mod.json", "README.md", "LICENSE", "src/a.mbt", "tests/a_test.mbt", "examples/demo.mbt", "benches/b.mbt"} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("Create zip entry: %v", err)
		}
		_, _ = w.Write([]byte("x"))
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("Close zip: %v", err)
	}

	summary, err := SummarizeZip(buf.Bytes())
	if err != nil {
		t.Fatalf("SummarizeZip: %v", err)
	}
	if !summary.MoonMod || !summary.Readme || !summary.License || summary.Sources == 0 || summary.Tests == 0 || summary.Examples == 0 || summary.Benches == 0 {
		t.Fatalf("summary = %#v", summary)
	}
}
