package badman_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/m-mizutani/badman"
	"github.com/m-mizutani/badman/source"
)

func Example() {
	man := badman.New()

	if err := man.Download(source.DefaultSet); err != nil {
		log.Fatal("Fail to download:", err)
	}

	fd, err := os.Open("ipaddrs_in_traffic_logs.txt")
	if err != nil {
		log.Fatal("Fail to open a file:", err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		entities, err := man.Lookup(scanner.Text())
		if err != nil {
			log.Fatal("Fail to lookup:", err)
		}

		if len(entities) > 0 {
			log.Printf("Matched %s in %s list (reason: %s)\n",
				entities[0].Name, entities[0].Src, entities[0].Reason)
		}
	}
}

func ExampleBadMan_Insert() {
	man := badman.New()

	if err := man.Insert(badman.BadEntity{
		Name:    "10.0.0.1",
		SavedAt: time.Now(),
		Src:     "It's me",
		Reason:  "testing",
	}); err != nil {
		log.Fatal("Fail to insert an entity:", err)
	}

	entities, err := man.Lookup("10.0.0.1")
	if err != nil {
		log.Fatal("Fail to lookup an entity:", err)
	}

	fmt.Println(entities[0].Name)
	// Output: 10.0.0.1
}

func ExampleBadMan_Dump() {
	//SetUp
	tmp, err := ioutil.TempFile("", "*.dat")
	if err != nil {
		log.Fatal(err)
	}
	tmp.Close()

	// Example
	man := badman.New()

	if err := man.Insert(badman.BadEntity{
		Name:    "orange.example.com",
		SavedAt: time.Now(),
		Src:     "clock",
	}); err != nil {
		log.Fatal("Fail to insert an entity:", err)
	}

	wfd, err := os.Create(tmp.Name())
	if err != nil {
		log.Fatal("Fail to create a file:", err)
	}

	// Save current repository to a file
	if err := man.Dump(wfd); err != nil {
		log.Fatal("Fail to dump repository")
	}
	wfd.Close()

	// Restore repository from a file
	rfd, err := os.Open(tmp.Name())
	if err != nil {
		log.Fatal("Fail to open a serialized data file:", err)
	}

	if err := man.Load(rfd); err != nil {
		log.Fatal("Fail to load repository")
	}

	entities, _ := man.Lookup("orange.example.com")

	fmt.Println(entities[0].Name)

	// TearDown
	rfd.Close()
	os.Remove(tmp.Name())

	// Output: orange.example.com
}
