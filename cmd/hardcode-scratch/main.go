package main

import (
	"context"
	"log"
	"os"

	"github.com/moby/buildkit/client/llb"
)

func main() {
	dt, err := createLLBState().Marshal(context.TODO(), llb.LinuxAmd64)
	if err != nil {
		log.Fatalf("Failed to marshal llb: %s", err)
	}

	llb.WriteTo(dt, os.Stdout)
}

func createLLBState() llb.State {
	return llb.Scratch().
		File(llb.Copy(llb.Local("context"), "README.md", "README.md"))
}
