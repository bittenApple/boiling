# boiling
An incremented id generator based on etcd, mainly relied on the behavior that etcd increments the version of key when any modification(put call) occured

[![Build Status](https://travis-ci.org/bittenApple/boiling.svg?branch=master)](https://travis-ci.org/bittenApple/boiling)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/bittenApple/boiling?status.svg)](https://godoc.org/github.com/bittenApple/boiling)
[![Go Report Card](https://goreportcard.com/badge/github.com/bittenApple/boiling)](https://goreportcard.com/report/github.com/bittenApple/boiling)

## Import
```
go get -u github.com/bittenApple/boiling
```

## Usage example
First you need a etcd instance running on http://localhost:2379

```
package main

import (
	"fmt"
	"log"

	"github.com/bittenApple/boiling"
)

func main() {
	o := &boiling.Options{
		Endpoints: []string{"http://localhost:2379"},
	}
	cli, err := boiling.NewClient(o)
	if err != nil {
		log.Printf("boiling client failed")
		return
	}
	fmt.Println(cli.GetId())
}
```
