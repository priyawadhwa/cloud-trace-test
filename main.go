package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
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

	_, flush, err := texporter.InstallNewPipeline(
		[]texporter.Option{
			texporter.WithProjectID("priya-wadhwa"),
		},
		// This example code uses sdktrace.AlwaysSample sampler to sample all traces.
		// In a production environment or high QPS setup please use ProbabilitySampler
		// set at the desired probability.
		// Example:
		// sdktrace.WithConfig(sdktrace.Config {
		//     DefaultSampler: sdktrace.ProbabilitySampler(0.0001),
		// })
		sdktrace.WithConfig(sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
		}),
		// other optional provider options
	)
	if err != nil {
		log.Fatalf("texporter.InstallNewPipeline: %v", err)
	}
	defer flush()

	fileName := "service/trace_context.json"
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "reading %s", fileName)
	}
	// marshal into propagator type
	tmc := tmc.New()
	if err := json.Unmarshal(contents, &tmc); err != nil {
		return errors.Wrap(err, "unmarshalling propagator")
	}
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
