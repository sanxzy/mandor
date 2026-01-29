package main

import (
	"os"
	"mandor/internal/cmd"
)

func main() {
	os.Exit(cmd.ExecuteWithCode())
}
