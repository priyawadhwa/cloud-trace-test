package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/priyawadhwa/cloud-census-test/tmc"
	"go.opentelemetry.io/otel/propagators"
)

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	fileName := "service1/baggage.json"
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "reading %s", fileName)
	}
	// marshal into propagator type
	propagator := propagators.TraceContext{}
	if err := json.Unmarshal(contents, &propagator); err != nil {
		return errors.Wrap(err, "unmarshalling propagator")
	}

	ctx := propagator.Extract(context.Background(), tmc.New())

	return nil
}
