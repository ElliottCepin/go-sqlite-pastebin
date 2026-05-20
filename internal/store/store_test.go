package store
import (
	"github.com/ElliottCepin/go-sqlite-pastebin/internal/server"
	"testing"
	"sync"
	"fmt"
)

type Storage struct {
	Name string
	New func() server.Store
}

type Case struct {
	Name string
	Inst func(*testing.T, server.Store)
}

func TestRoundTrip(t *testing.T) {
	stores := make([]Storage, 0, 0)	
	tests := make([]Case, 0, 0)	
	ms := Storage{
		Name: "MemoryStore",
		New: func() server.Store {
			return NewMemoryStore()
		},
	}
	
	rt := Case {
		Name: "RoundTrip",
		Inst: subtestRoundTrip,
	}

	del := Case {
		Name: "Delete",
		Inst: subtestDelete,
	}

	col := Case {
		Name: "Collision",
		Inst: subtestCollision,
	}

	co := Case {
		Name: "ConcurrentOps",
		Inst: subtestConcurrentOps,
	}

	stores = append(stores, ms)
	tests = append(tests, rt)	
	tests = append(tests, del)
	tests = append(tests, col)
	tests = append(tests, co)

	for _, impl := range stores {
		for _, test := range tests {
			t.Run("Test" + test.Name + ":" + impl.Name, func (t *testing.T) { test.Inst(t, impl.New()) } )
		}
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

func subtestDelete(t *testing.T, s server.Store) {
	in := "abcd" 
	out := "efgh"

	if err := s.Save(t.Context(), in, out); err != nil {
		t.Errorf("An unexpected error occured during s.Save: %v", err)
	}

	if err := s.Delete(t.Context(), in); err != nil {
		t.Errorf("An unexpected error occured during s.Delete: %v", err)
	}

	val, err := s.Get(t.Context(), in)
	if (err == nil) {
		t.Errorf("Expected an error, got %v instead", val)
	}
	
}

func subtestCollision (t *testing.T, s server.Store) {
	in := "abcd" 
	out1 := "efgh"
	out2 := "ijkl"

	if err := s.Save(t.Context(), in, out1); err != nil {
		t.Errorf("An unexpected error occured during s.Save: %v", err)
	}

	if err := s.Save(t.Context(), in, out2); err == nil {
		t.Errorf("Expected an error, but got none")
	}

	if val, _ := s.Get(t.Context(), in); val == out2 {
		t.Errorf("Initial value %v got overwritten to %v", out1, out2)
	}
	
}

func subtestConcurrentOps(t *testing.T, s server.Store) {
	var wg sync.WaitGroup
	for i := 0; i<100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			slug := fmt.Sprintf("slug-%v", i % 10)
			_ = s.Save(t.Context(), slug, "content")
			_, _ = s.Get(t.Context(), slug)
			_ = s.Delete(t.Context(), slug)
		}(i)
	}
	wg.Wait()
}
