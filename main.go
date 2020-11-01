package main

import (
	"fmt"
)

func main() {
	types := []string{"plain", "base64", "sip008", "clash"}
	for _, v := range types {
		fmt.Printf("==== MODE:%s ====\n%s\n", v, GenerateSubscription(v))
	}
}
