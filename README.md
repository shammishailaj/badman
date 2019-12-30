# badman [![Travis-CI](https://travis-ci.org/m-mizutani/badman.svg)](https://travis-ci.org/m-mizutani/badman) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/badman)](https://goreportcard.com/report/github.com/m-mizutani/badman) [![GoDoc](https://godoc.org/github.com/m-mizutani/badman?status.svg)](https://godoc.org/github.com/m-mizutani/badman)


**B**lacklisted **A**ddress and **D**omain name **Man**ager is tool to manage blacklist network entities. The tool provides downloader/save 

## Library Examples

### Getting Started

```go
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/m-mizutani/badman"
	"github.com/m-mizutani/badman/source"
)

func main() {
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
		entity, err := man.Lookup(scanner.Text())
		if err != nil {
			log.Fatal("Fail to lookup:", err)
		}

		if entity != nil {
			log.Printf("Matched %s in %s list (reason: %s)\n",
				entity.Name, entity.Src, entity.Reason)
		}
	}
}
```

`ipaddrs_in_traffic_logs.txt` includes only IP addresses line by line. `BadMan` downloads backlists from default sources (blacklist providers) and store entities in the blacklist into own repository. Default settings are below. After downloading blacklist, `Lookup` method is enabled to search given name (IP address or domain name, both are accepted) from the repository.

- Default sources: `MVPS`, `MalwareDomains`, `URLhausRecent`, `URLhausOnline`
- Default repository: `inMemoryRepository`

### Insert a new bad entity one by one

```go
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

	// Output:
	// 10.0.0.1
	fmt.Println(entities[0].Name)
```

### Save and Restore

```go
	man := badman.New()

	// Save current repository to a file
	if err := man.Dump(wfd); err != nil {
		log.Fatal("Fail to dump repository")
	}
	wfd.Close()

	// Restore repository from a file
	rfd, err := os.Open("repo.dat")
	if err := man.Load(rfd); err != nil {
		log.Fatal("Fail to load repository")
    }
```

## Use case

### Stateless (Serveless model)

![Serverless Architecture for AWS](https://user-images.githubusercontent.com/605953/71566177-b844e400-2af8-11ea-8c65-bc5e8757be9e.png)

### Stateful (Server model)

## Terms of Use for Data Sources

The tool uses several online blacklist sites. They have each own Terms of Use and please note you need to understand their policy before operating in your environment. A part of Terms of Use regarding usage policy is following.

### Winhelp2002 ( `MVPS` )

> Disclaimer: this file is free to use for personal use only. Furthermore it is NOT permitted to copy any of the contents or host on any other site without permission ormeeting the full criteria of the below license terms.
>
> This work is licensed under the Creative Commons Attribution-NonCommercial-ShareAlike License.
> https://creativecommons.org/licenses/by-nc-sa/4.0/

http://winhelp2002.mvps.org/hosts.txt

### DNS-BH – Malware Domain Blocklist by RiskAnalytics ( `MalwareDomains` )

> This malware block lists provided here are for free for noncommercial use as part of the fight against malware.
>
> Any use of this list commercially is strictly prohibited without prior approval. (It’s OK to use this list on an internal DNS server for which you are not charging).

http://www.malwaredomains.com/?page_id=1508

### URLhaus ( `URLhausRecent`, `URLhausOnline` )

> All datasets offered by URLhaus can be used for both, commercial and non-commercial purpose without any limitations (CC0)

https://urlhaus.abuse.ch/api/#tos
