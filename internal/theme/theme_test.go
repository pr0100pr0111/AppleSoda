package theme

import (
	"testing"
)

func TestGet(t *testing.T) {
	th := Get("catppuccin")
	if th.Name != "catppuccin" {
		t.Fatalf("expected catppuccin, got %s", th.Name)
	}

	fallback := Get("unknown")
	if fallback.Name != "apple" {
		t.Fatalf("expected fallback apple, got %s", fallback.Name)
	}
}

func TestNextTheme(t *testing.T) {
	first := Names()[0]
	next := NextTheme(first)
	if next == first {
		t.Fatalf("expected next theme to be different, got %s", next)
	}

	last := Names()[len(Names())-1]
	cycled := NextTheme(last)
	if cycled != first {
		t.Fatalf("expected cycle to %s, got %s", first, cycled)
	}
}
