package parser

import "testing"

func TestParseAsync(t *testing.T) {
    // CWD = dredger/parser
    doc, err := ParseAsyncAPISpecFile("../examples/hello_async.yaml")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if doc.Info.Title != "HelloEvent" {
        t.Fatalf("got title %q, want %q", doc.Info.Title, "HelloEvent")
    }
}
