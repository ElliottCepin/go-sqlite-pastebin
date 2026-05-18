package store
import (
	"github.com/ElliottCepin/go-sqlite-pastebin/internal/server"
	"testing"
)

type Case struct {
	Name string
	New func() server.Store
}

func TestRoundTrip(t *testing.T) {
	tests := make([]Case, 0, 0)	
	ms := Case{
		Name: "MemoryStore",
		New: func() server.Store {
			return NewMemoryStore()
		},
	}
	tests = append(tests, ms)

	for _, impl := range tests {
		t.Run(impl.Name, func (t *testing.T) { subtestRoundTrip(t, impl.New()) } )
	}
}

func subtestRoundTrip(t *testing.T, s server.Store) {
	in := "abcd"
	out := "efgh"
	if err := s.Save(t.Context(), in, out); err != nil {
		t.Errorf("An unexpected error occured during s.Save: %v", err)
	}
	val, err := s.Get(t.Context(), in)
	if (err != nil) {
		t.Errorf("An unexpected error occured during s.Get: %v", err)
	}

	if (val != out) {
		t.Errorf("s.Get returned the wrong value. Expected %v, got %v", out, val)
	}
}
