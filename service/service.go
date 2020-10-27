package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/pkg/errors"
	"github.com/priyawadhwa/cloud-trace-test/tmc"
)

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	// Create exporter and trace provider pipeline, and register provider.
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

	t := global.Tracer("container-tools")
	ctx, span := t.Start(context.Background(), "first_service")

	// we want to propagate this trace ID to another tool so save it to a file
	propagator := propagators.TraceContext{}
	tmc := tmc.New()
	propagator.Inject(ctx, tmc)
	// save in a file
	fileName := "trace_context.json"
	contents, err := json.Marshal(tmc)
	if err != nil {
		return errors.Wrap(err, "marshalling propagator")
	}
	if err := ioutil.WriteFile(fileName, contents, 0644); err != nil {
		return errors.Wrap(err, "writing file")
	}

	fmt.Println("trace id:", span.SpanContext().TraceID)
	fmt.Println("sleeping 20 seconds...")
	time.Sleep(20 * time.Second)
	span.End()
	return nil
}
