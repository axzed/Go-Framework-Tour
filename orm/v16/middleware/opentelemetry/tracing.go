//go:build v16
package opentelemetry

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const defaultInstrumentationName = "gitee.com/geektime-geekbang/geektime-go/orm/middleware/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (b *MiddlewareBuilder) Build() orm.Middleware {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			tbl := qc.Model.TableName
			reqCtx, span := b.Tracer.Start(ctx, qc.Type+ "-" + tbl, trace.WithAttributes())
			defer span.End()
			span.SetAttributes(attribute.String("component", "orm"))
			q, err := qc.builder.Build()
			if err != nil {
				span.RecordError(err)
			}
			span.SetAttributes(attribute.String("table", tbl))
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}
			return next(reqCtx, qc)
		}
	}
}