# badman [![Travis-CI](https://travis-ci.org/m-mizutani/badman.svg)](https://travis-ci.org/m-mizutani/badman) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/badman)](https://goreportcard.com/report/github.com/m-mizutani/badman) [![GoDoc](https://godoc.org/github.com/m-mizutani/badman?status.svg)](https://godoc.org/github.com/m-mizutani/badman)


**B**lacklisted **A**ddress and **D**omain name **Man**ager is tool to manage blacklist network entities. The tool provides download, save and restore capability about blacklist of IP address and domain name.

## Examples

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
```

`ipaddrs_in_traffic_logs.txt` includes only IP addresses line by line. In this case, `BadMan` downloads backlists prepared in this package and store entities (IP address or domain name) in the blacklist into own repository. After downloading blacklist, `Lookup` method is enabled to search given name (IP address or domain name, both are accepted) from the repository.

Default settings are following.

- Default sources: `MVPS`, `MalwareDomains`, `URLhausRecent` and `URLhausOnline`
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

	wfd, err := os.Create("repo.dat")
	if err != nil {
		log.Fatal("Fail to create a file:", err)
	}

	// Save current repository to a file
	if err := man.Dump(wfd); err != nil {
		log.Fatal("Fail to dump repository")
	}
	wfd.Close()

	// Restore repository from a file
	rfd, err := os.Open("repo.dat")
	if err != nil {
		log.Fatal("Fail to open a serialized data file:", err)
	}

	if err := man.Load(rfd); err != nil {
		log.Fatal("Fail to load repository")
	}
```

### Change blacklist sources

```go
	man := badman.New()
	set := []badman.Source{
		source.NewURLhausRecent(),
		source.NewURLhausOnline(),
	}
	if err := man.Download(set); err != nil {
		log.Fatal("Fail to download:", err)
	}
```

Usually you can use `source.DefaultSet` to download all badman supported blacklist providers (sources). However, if you want to use specific sources, you can choose your preffered sources. For example, above sample code downloads only URLhaus blacklist.

### Change repository

```go
	dynamoTableRegion, dynamoTableName := "ap-northeast-1", "your-table-name"
	man := badman.New()
	man.ReplaceRepository(badman.NewDynamoRepository(dynamoTableRegion, dynamoTableName))
```

Repository can be replaced by `ReplaceRepository()` with other repository. When replacing, blacklist data in old repository is NOT copied to new repository. Below 2 repositories are prepared in this pacakge.

- `inMemoryRepository`
- `dynamoRepository`

Also, you can use own repository that is implemented `badman.Repository` interface.

## Use case

Basically `badman` should be used as library and a user need to implement own program by leveraging `badman`.

### Stateless (Serveless model)

![Serverless Architecture for AWS](https://user-images.githubusercontent.com/605953/71566177-b844e400-2af8-11ea-8c65-bc5e8757be9e.png)

In this case, there are 2 Lambda function. 1st (left side) function retrieves blacklist and saves dumped blacklist data to S3 periodically. 2nd (right side) Lambda function is invoked by S3 ObjectCreated event of traffic log file. After that, the Lambda function downloads both of dumped blacklist data and log file and check if the IP addresses of traffic logs exist in blacklist. If existing, lambda notify it to an administrator via communication tool, such as Slack.

### Stateful (Server model)

<img width="640" alt="Server based Architecture for AWS" src="https://user-images.githubusercontent.com/605953/71639479-b3c82900-2cba-11ea-9fb9-08201edf9271.png">

A major advantage of server model is stream processing for real-time capability. Above serverless model has latency because buffering is required to assemble tiny log data to one object before uploading to S3. Generally, using the above model, the delay will be on the order of minutes. Therefore, it is recommended to use the server model when lower-latency processing is required.

This program has a fluentd interface and receives traffic logs via fluentd. After that, use badman to check the traffic log for IP addresses included in the blacklist. Blacklist expects to be updated periodically, and uses DynamoDB as the repository so that it can recover even if the host running the program (in this case, EC2) crashes.


## Terms of Use for Data Sources

The tool uses several online blacklist sites. They have each own Terms of Use and please note you need to understand their policy before operating in your environment. A part of their Terms of Use regarding usage policy is below.

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


## License

- [MIT License](./LICENSE)
- Author: Masayoshi Mizutani < mizutani@sfc.wide.ad.jp >