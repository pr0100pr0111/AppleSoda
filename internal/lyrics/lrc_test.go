package lyrics

import (
	"testing"
)

func TestParseLRCBasic(t *testing.T) {
	src := "[00:10.00]Hello world\n[00:20.50]Second line"
	lines := ParseLRC(src)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].TimeSec != 10.0 {
		t.Errorf("expected 10.0, got %f", lines[0].TimeSec)
	}
	if lines[0].Text != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", lines[0].Text)
	}
	if lines[1].TimeSec != 20.5 {
		t.Errorf("expected 20.5, got %f", lines[1].TimeSec)
	}
}

func TestParseLRCMultipleTimestamps(t *testing.T) {
	src := "[00:10.00][00:20.00]chorus line"
	lines := ParseLRC(src)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].TimeSec != 10.0 {
		t.Errorf("expected 10.0, got %f", lines[0].TimeSec)
	}
	if lines[1].TimeSec != 20.0 {
		t.Errorf("expected 20.0, got %f", lines[1].TimeSec)
	}
	if lines[0].Text != "chorus line" || lines[1].Text != "chorus line" {
		t.Errorf("unexpected text: %q %q", lines[0].Text, lines[1].Text)
	}
}

func TestParseLRCOffset(t *testing.T) {
	src := "[offset:1000]\n[00:10.00]line"
	lines := ParseLRC(src)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0].TimeSec != 11.0 {
		t.Errorf("expected 11.0 (offset 1000ms), got %f", lines[0].TimeSec)
	}
}

func TestParseLRCNegativeOffset(t *testing.T) {
	src := "[offset:-500]\n[00:10.00]line"
	lines := ParseLRC(src)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0].TimeSec != 9.5 {
		t.Errorf("expected 9.5 (offset -500ms), got %f", lines[0].TimeSec)
	}
}

func TestParseLRCEmptyTextLines(t *testing.T) {
	src := "[00:10.00]\n[00:20.00]line"
	lines := ParseLRC(src)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (empty text kept), got %d", len(lines))
	}
	if lines[0].Text != "" {
		t.Errorf("expected empty text, got %q", lines[0].Text)
	}
}

func TestParseLRCNonTimestampLinesIgnored(t *testing.T) {
	src := "This is a plain line\n[00:10.00]timestamped line\nAnother plain line"
	lines := ParseLRC(src)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestParseLRCSorted(t *testing.T) {
	src := "[00:30.00]third\n[00:10.00]first\n[00:20.00]second"
	lines := ParseLRC(src)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].TimeSec != 10.0 || lines[1].TimeSec != 20.0 || lines[2].TimeSec != 30.0 {
		t.Errorf("lines not sorted: %v", lines)
	}
}

func TestActiveIndex(t *testing.T) {
	lines := []Line{
		{TimeSec: 0, Text: "intro"},
		{TimeSec: 10, Text: "verse"},
		{TimeSec: 20, Text: "chorus"},
	}

	if ActiveIndex(lines, -1) != 0 {
		t.Error("expected 0 for negative pos")
	}
	if ActiveIndex(lines, 0) != 0 {
		t.Error("expected 0 at t=0")
	}
	if ActiveIndex(lines, 5) != 0 {
		t.Error("expected 0 at t=5")
	}
	if ActiveIndex(lines, 10) != 1 {
		t.Error("expected 1 at t=10")
	}
	if ActiveIndex(lines, 15) != 1 {
		t.Error("expected 1 at t=15")
	}
	if ActiveIndex(lines, 20) != 2 {
		t.Error("expected 2 at t=20")
	}
	if ActiveIndex(lines, 999) != 2 {
		t.Error("expected 2 at t=999")
	}
}

func TestActiveIndexEmpty(t *testing.T) {
	if ActiveIndex(nil, 10) != 0 {
		t.Error("expected 0 for nil lines")
	}
	if ActiveIndex([]Line{}, 10) != 0 {
		t.Error("expected 0 for empty lines")
	}
}
