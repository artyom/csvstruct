package csvstruct

import (
	"testing"
	"time"
)

func TestDecoder(t *testing.T) {
	var dst testDst
	decoder, err := NewDecoder([]string{"name", "age"}, &dst)
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

type testDstValue struct {
	Name string `csv:"name"`
	Time myTime `csv:"time"`
}

type myTime time.Time

func (t *myTime) Set(s string) error {
	t2, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	*t = myTime(t2)
	return nil
}

func TestValue(t *testing.T) {
	var dst testDstValue
	decoder, err := NewDecoder([]string{"name", "time"}, &dst)
	if err != nil {
		t.Fatal(err)
	}
	if dst != (testDstValue{}) {
		t.Fatalf("NewDecoder modified destination: %+v", dst)
	}
	if err := decoder([]string{"boom", "2020-04-01T12:00:00Z"}, &dst); err != nil {
		t.Fatal(err)
	}
	timestamp := time.Time(dst.Time)
	if ref := time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC); !timestamp.Equal(ref) {
		t.Fatalf("dst.Time %v differs from reference time %v)", timestamp, ref)
	}
}

func BenchmarkDecoder(b *testing.B) {
	decoder, err := NewDecoder([]string{"name", "age"}, &sink)
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
