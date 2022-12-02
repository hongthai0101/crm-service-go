package utils

import "fmt"

func Debug(a ...any) {
	println("\n ================================================================== \n")
	for i := 0; i < len(a); i++ {
		fmt.Printf("%v %T \n", a[i], a[i])
	}
	println("\n ================================================================== \n")
}
