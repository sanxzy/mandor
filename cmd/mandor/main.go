package main

import (
	"mandor/internal/cmd"
	"os"
)

func main() {
	os.Exit(cmd.ExecuteWithCode())
}
