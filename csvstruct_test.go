package csvstruct

import (
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	var dst testDst
	scan, err := NewScanner([]string{"name", "age"}, &dst)
	if err != nil {
		t.Fatal(err)
	}
	if dst != (testDst{}) {
		t.Fatalf("NewScanner modified destination: %+v", dst)
	}
	if err := scan([]string{"John", "42"}, &dst); err != nil {
		t.Fatal(err)
	}
	want := testDst{Name: "John", Age: 42}
	if dst != want {
		t.Fatalf("got: %+v, want: %+v", dst, want)
	}
}

func BenchmarkScanner(b *testing.B) {
	scan, err := NewScanner([]string{"name", "age"}, &sink)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := scan([]string{"John", "42"}, &sink); err != nil {
			b.Fatal(err)
		}
	}
}

var sink testDst

type testDst struct {
	Name   string `csv:"name"`
	Age    uint8  `csv:"age"`
	Ignore bool
}

func TestScanner_DifferentType(t *testing.T) {
	scan, err := NewScanner([]string{"name", "age"}, &testDst{})
	if err != nil {
		t.Fatal(err)
	}
	type testDst2 struct {
		Name   string `csv:"name"`
		Age    uint8  `csv:"age"`
		Ignore bool
		_      byte // <- differs from testDst type
	}
	var dst testDst2
	defer func() {
		p := recover()
		if p == nil {
			t.Fatal("using different types in NewScanner and Scanner should have panicked, but it's not")
		}
		if s, ok := p.(string); !ok || !strings.Contains(s, "different type") {
			t.Fatalf("unexpected panic value: %v", p)
		}
	}()
	t.Fatalf("Scanner returned on different types: %v", scan([]string{"John", "42"}, &dst))
}
