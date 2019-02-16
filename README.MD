# A simple pubsub for Go

``` Go
package main

import (
	"fmt"
	"github.com/StevenZack/pubsub"
)

func main() {
	go handle(0)
	go handle(1)
	go handle(2)
	s := ""
	for {
		fmt.Scanf("%s", &s)
		pubsub.Pub("chanId", s)
		if s == "quit" {
			break
		}
	}
}
func handle(i int) {
	ps := pubsub.NewPubSub()
	ps.Sub("chanId", func(v interface{}) {
		if s, ok := v.(string); ok {
			println(i, ":", s)
			if s == "quit" {
				ps.UnSub()
			}
		} else {
			println(i, ": type error")
		}
	})
}



```

