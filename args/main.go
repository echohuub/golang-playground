package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) > 0 {
		for index, arg := range os.Args {
			fmt.Printf("args[%d]=%v\n", index, arg)
		}
	}
}
