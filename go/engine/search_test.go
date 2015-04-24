package engine

import (
	"testing"

	keybase_1 "github.com/keybase/client/protocol/go"
)

func TestSearch(t *testing.T) {
	tc := SetupEngineTest(t, "btc")
	defer tc.Cleanup()

	ctx := &Context{
		LogUI: G.UI.GetLogUI(),
	}
	// This twitter handle is used by t_alice and t_charlie.
	e := NewSearchEngine("tacovontaco")
	err := RunEngine(e, ctx)
	if err != nil {
		t.Fatal(err)
	}
	results := e.GetResults()

	if len(results) < 2 {
		// The test DB could contain other random test users, namely Max's ;)
		t.Fatalf("Expected at least 2 search results for 'tacovontaco'. Got %d.", len(results))
	}

	var t_alice keybase_1.UserSummary
	var t_charlie keybase_1.UserSummary
	for _, summary := range results {
		if summary.Username == "t_alice" {
			t_alice = summary
		} else if summary.Username == "t_charlie" {
			t_charlie = summary
		}
	}

	if len(t_alice.Proofs.Social) != 2 {
		t.Fatalf("Expected 2 proofs for t_alice, got %d.", len(t_alice.Proofs.Social))
	}
	if len(t_charlie.Proofs.Social) != 2 {
		t.Fatalf("Expected 2 proofs for t_alice, got %d.", len(t_charlie.Proofs.Social))
	}
}
