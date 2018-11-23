# boiling
An incremented id generator based on etcd, mainly relied on the behavior that etcd increments the version of key when any modification(put call) occured

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
