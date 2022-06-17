package Service

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type metricsMiddlewareServer struct {
	next      Service
	counter   metrics.Counter
	histogram metrics.Histogram
}

func (m metricsMiddlewareServer) GetAllWeightRecord(ctx context.Context) (out AllDocumentsAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetAllWeightRecord"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetAllWeightRecord(ctx)
	return
}

func (m metricsMiddlewareServer) GetWeightRecordPage(ctx context.Context, in Page) (out AllDocumentsPageAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetWeightRecordPage"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetWeightRecordPage(ctx, in)
	return
}

func (m metricsMiddlewareServer) GetParameter(ctx context.Context) (out AllParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetParameter"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetParameter(ctx)
	return
}

func (m metricsMiddlewareServer) AddNewRecord(ctx context.Context, in NewRecord) (out NewRecordAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddNewRecord"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.AddNewRecord(ctx, in)
	return
}

func NewMetricsMiddlewareServer(counter metrics.Counter, histogram metrics.Histogram) NewMiddlewareServer {
	return func(service Service) Service {
		return metricsMiddlewareServer{
			next:      service,
			counter:   counter,
			histogram: histogram,
		}
	}
}
