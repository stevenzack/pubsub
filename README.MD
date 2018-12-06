# A simple pubsub for Go

``` Go
package main

import (
	"fmt"

	"github.com/StevenZack/pubsub"
)

func main() {
	go ls(0)
	go ls(1)
	go ls(2)
	s := ""
	for {
		fmt.Scanf("%s", &s)
		pubsub.Pub("chan", s)
	}
}
func ls(i int) {
	chanId := "chan"
	ps := pubsub.NewPubSub()
	defer ps.UnSub(chanId)
	ps.Sub(chanId, func(b []byte) {
		fmt.Println(i, ":", string(b))
	})
}


```

