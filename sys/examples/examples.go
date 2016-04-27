package main

import (
	"fmt"

	"github.com/colinyl/ars/sys"
)

func main() {
	fmt.Println(sys.GetMemory())
	fmt.Println(sys.GetCPU())
	fmt.Println(sys.GetDisk())
}
