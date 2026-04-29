package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	d := Diff(base, target)

	if len(d.Added) != 0 || len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v", d.Added, d.Removed, d.Changed)
	}
	if len(d.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(d.Unchanged))
	}
}

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "newval"}

	d := Diff(base, target)

	if len(d.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(d.Added))
	}
	if d.Added["NEW_KEY"] != "newval" {
		t.Errorf("unexpected added value: %s", d.Added["NEW_KEY"])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "oldval"}
	target := map[string]string{"FOO": "bar"}

	d := Diff(base, target)

	if len(d.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(d.Removed))
	}
	if d.Removed["OLD_KEY"] != "oldval" {
		t.Errorf("unexpected removed value: %s", d.Removed["OLD_KEY"])
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "baz"}

	d := Diff(base, target)

	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(d.Changed))
	}
	pair := d.Changed["FOO"]
	if pair[0] != "bar" || pair[1] != "baz" {
		t.Errorf("unexpected change: %v", pair)
	}
}

func TestWriteDiff_NoDiff(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	var buf bytes.Buffer
	WriteDiff(&buf, Diff(base, base))
	if !strings.Contains(buf.String(), "(no differences)") {
		t.Errorf("expected no-differences message, got: %s", buf.String())
	}
}

func TestWriteDiff_ShowsChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "99", "C": "3"}
	var buf bytes.Buffer
	WriteDiff(&buf, Diff(base, target))
	out := buf.String()
	if !strings.Contains(out, "+ C=3") {
		t.Errorf("expected added C, got: %s", out)
	}
	if !strings.Contains(out, "- B=2") {
		t.Errorf("expected removed B, got: %s", out)
	}
	if !strings.Contains(out, "~ A: 1 -> 99") {
		t.Errorf("expected changed A, got: %s", out)
	}
}
