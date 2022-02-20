package main

import (
	"github.com/junaidrahim/argodocs/workflow"
)

func main() {
	_, err := workflow.ParseFiles("/tmp/workflows/*.yaml")
	if err != nil {
		return
	}
}
