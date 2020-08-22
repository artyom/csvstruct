package csvstruct

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

type Person struct {
	FirstName string `csv:"first_name"`
	LastName  string `csv:"last_name"`
}

const Table = `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`

func ExampleScanner() {
	r := csv.NewReader(strings.NewReader(Table))
	header, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	scan, err := NewScanner(header, &Person{})
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
		var person Person
		if err := scan(record, &person); err != nil {
			log.Fatal(err)
		}
		fmt.Println(person.FirstName, person.LastName)
	}
	// Output:
	//
	// Rob Pike
	// Ken Thompson
	// Robert Griesemer
}
