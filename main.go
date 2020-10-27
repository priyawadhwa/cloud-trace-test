package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/priyawadhwa/cloud-trace-test/tmc"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/propagators"
)

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	fileName := "service/baggage.json"
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "reading %s", fileName)
	}
	// marshal into propagator type
	tmc := tmc.New()
	if err := json.Unmarshal(contents, &tmc); err != nil {
		return errors.Wrap(err, "unmarshalling propagator")
	}
	fmt.Println(tmc)
	propagator := propagators.TraceContext{}

	ctx := propagator.Extract(context.Background(), tmc)
	t := global.Tracer("container-tools")
	_, span := t.Start(ctx, "second_tool")
	defer span.End()
	fmt.Println("sleeping 10 seconds...")
	fmt.Println(span.SpanContext().TraceID)
	time.Sleep(10 * time.Second)
	return nil
}
