package csvstruct

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Event struct {
	Name string `csv:"name"`
	Time myTime `csv:"timestamp"`
}

type myTime struct {
	time.Time
}

func (t *myTime) Set(s string) error {
	t2, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = t2
	return nil
}

const EventsTable = `name,timestamp
login,2020-08-23T10:57:27Z
logout,2020-08-23T11:07:09Z
`

func ExampleValue() {
	r := csv.NewReader(strings.NewReader(EventsTable))
	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	scan, err := NewScanner(header, &Event{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var event Event
		if err := scan(record, &event); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", event.Name, event.Time.Format("2006-01-02 15:04:05 MST"))
	}
	// Output:
	//
	// login: 2020-08-23 10:57:27 UTC
	// logout: 2020-08-23 11:07:09 UTC
}
