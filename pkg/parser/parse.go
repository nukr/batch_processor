package parser

import "fmt"

// Parse ...
func Parse() {
	operator := "$push"
	switch operator {
	case "$push":
		fmt.Println("do push")
	case "$pull":
		fmt.Println("do pull")
	default:
		fmt.Println("error")
	}
}
