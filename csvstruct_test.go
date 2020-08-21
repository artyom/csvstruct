package csvstruct

import "testing"

func TestDecoder(t *testing.T) {
	var dst testDst
	decoder, err := NewDecoder([]string{"name", "age"}, dst)
	if err != nil {
		t.Fatal(err)
	}
	if dst != (testDst{}) {
		t.Fatalf("NewDecoder modified destination: %+v", dst)
	}
	if err := decoder([]string{"John", "42"}, &dst); err != nil {
		t.Fatal(err)
	}
	want := testDst{Name: "John", Age: 42}
	if dst != want {
		t.Fatalf("got: %+v, want: %+v", dst, want)
	}
}

func BenchmarkDecoder(b *testing.B) {
	decoder, err := NewDecoder([]string{"name", "age"}, sink)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := decoder([]string{"John", "42"}, &sink); err != nil {
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
