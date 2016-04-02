package main

import "fmt"

// START OMIT
const tmplString = `package main

import "fmt"

func main() {
	strs := %#v
	for _, str := range strs {
		fmt.Println(str)
	}
}
`

func main() {
	strs := make([]string, 2)
	strs[0] = "Hello"
	strs[1] = "world!"
	fmt.Printf(tmplString, strs) // HL
}

// END OMIT
