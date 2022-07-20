package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"k8s.io/klog/v2"

	mock "github.com/flyer103/buildkit-demo/pkg/frontend-mockerfile"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
)

var graph bool
var filename string

func main() {
	flag.BoolVar(&graph, "graph", false, "output a graph and exit")
	flag.StringVar(&filename, "filename", "Mockerfile.yaml", "the file to read from")
	flag.Parse()

	if graph {
		err := printGraph(filename, os.Stdout)
		if err != nil {
			klog.ErrorS(err, "Failed to print graph")
			os.Exit(1)
		}

		os.Exit(0)
	}

	err := grpcclient.RunFromEnvironment(appcontext.Context(), mock.Build)
	if err != nil {
		klog.ErrorS(err, "Failed to build")
		os.Exit(1)
	}

	klog.InfoS("Successfully build.")
}

func printGraph(filename string, out io.Writer) error {
	c, err := mock.NewFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	st, _ := mock.Mockerfile2LLB(c)
	dt, err := st.Marshal(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to marshal LLB: %s", err)
	}

	return llb.WriteTo(dt, out)
}
