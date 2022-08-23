package transport

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("occult.work/apt/transport")

func setSpanRequest(span trace.Span, request *Request) {
	span.SetAttributes(
		attribute.Stringer("request.source", request.Source),
		attribute.String("request.modified", request.Modified.Format(time.RFC1123)),
		attribute.String("request.target", request.Target),
	)
}
