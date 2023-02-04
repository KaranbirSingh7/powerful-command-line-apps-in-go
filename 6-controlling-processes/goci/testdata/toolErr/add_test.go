package add

import "testing"

func TestAdd(t *testing.T) {
	a := 3
	b := 5

	want := 8
	got := add(a, b)

	if got != want {
		t.Errorf("Got %q, want %q instead", got, want)
	}
}
