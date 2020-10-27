package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/pkg/errors"
	"github.com/priyawadhwa/cloud-census-test/tmc"
	"go.opentelemetry.io/otel/api/global"
)

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	exporter, err := texporter.NewExporter(texporter.WithProjectID("priya-wadhwa"))
	if err != nil {
		return errors.Wrap(err, "getting exporter")
	}
	tp, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))
	if err != nil {
		return errors.Wrap(err, "new provider")
	}
	tp.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
	global.SetTraceProvider(tp)
	t := global.TraceProvider().Tracer("container-tools")
	ctx, span := t.Start(context.Background(), "minikube_start")
	fmt.Println(span.SpanContext().TraceID)

	// do something with the baggage exporter
	propagator := propagators.TraceContext{}
	propagator.Inject(ctx, tmc.New())

	// save in a file
	fileName := "baggage.json"
	contents, err := json.Marshal(propagator)
	if err != nil {
		return errors.Wrap(err, "marshalling propagator")
	}
	if err := ioutil.WriteFile(fileName, contents, 0644); err != nil {
		return errors.Wrap(err, "writing file")
	}

	time.Sleep(10 * time.Minute)
	span.End()
	return nil
}
