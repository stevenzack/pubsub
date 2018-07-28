# A simple pubsub for Go

``` Go
package main

import (
	"fmt"
	"github.com/StevenZack/pubsub"
)

var ps = pubsub.NewPubSub()

func main() {
	go name(2)
	go name(3)
	go name(4)
	s := ""
	for {
		fmt.Scanf("%s", &s)
		if s == "quit" {
			break
		}
		ps.Pub(s)
	}
	fmt.Println("waiting for exit signal")
	fmt.Scanf("%s", &s)
}
func name(id int) {
	counter := 0
	chanId := pubsub.NewToken()
	ps.Sub(func(data interface{}) {
		counter++
		if counter == id {
			ps.UnSub(chanId)
			return
		}
		fmt.Println(id, counter, data, chanId)
	}, chanId)
}

```

