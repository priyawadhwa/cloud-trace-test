package trace

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/sdk/trace"
)

func NewProvider() (trace.Provider, error) {
	exporter, err := texporter.NewExporter(texporter.WithProjectID("priya-wadhwa"))
	if err != nil {
		return nil, errors.Wrap(err, "getting exporter")
	}
	tp, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))
	if err != nil {
		return nil, errors.Wrap(err, "new provider")
	}
	tp.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
	global.SetTraceProvider(tp)
	t := global.TraceProvider().Tracer("container-tools")
	return t, nil
}
