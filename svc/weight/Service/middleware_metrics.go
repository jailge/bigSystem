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

func (m metricsMiddlewareServer) UpdateTexture(ctx context.Context, id string, in Texture) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UpdateTexture"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.UpdateTexture(ctx, id, in)
	return
}

func (m metricsMiddlewareServer) UpdateProcess(ctx context.Context, id string, in Process) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UpdateProcess"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.UpdateProcess(ctx, id, in)
	return
}

func (m metricsMiddlewareServer) UpdatePurchaseStatus(ctx context.Context, id string, in PurchaseStatus) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UpdatePurchaseStatus"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.UpdatePurchaseStatus(ctx, id, in)
	return
}

func (m metricsMiddlewareServer) UpdateCraft(ctx context.Context, id string, in Craft) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "UpdateCraft"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.UpdateCraft(ctx, id, in)
	return
}

func (m metricsMiddlewareServer) GetAllPurchaseStatus(ctx context.Context) (out AllPurchaseStatusAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetAllPurchaseStatus"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetAllPurchaseStatus(ctx)
	return
}

func (m metricsMiddlewareServer) GetAllTexture(ctx context.Context) (out AllTextureAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetAllTexture"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetAllTexture(ctx)
	return
}

func (m metricsMiddlewareServer) GetAllProcess(ctx context.Context) (out AllProcessAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetAllProcess"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetAllProcess(ctx)
	return
}

func (m metricsMiddlewareServer) GetAllCraft(ctx context.Context) (out AllCraftAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "GetAllCraft"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.GetAllCraft(ctx)
	return
}

func (m metricsMiddlewareServer) DeletePurchaseStatusWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "DeletePurchaseStatusWithId"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.DeletePurchaseStatusWithId(ctx, in)
	return
}

func (m metricsMiddlewareServer) DeleteTextureWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "DeleteTextureWithId"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.DeleteTextureWithId(ctx, in)
	return
}

func (m metricsMiddlewareServer) DeleteProcessWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "DeleteProcessWithId"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.DeleteProcessWithId(ctx, in)
	return
}

func (m metricsMiddlewareServer) DeleteCraftWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "DeleteCraftWithId"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.DeleteCraftWithId(ctx, in)
	return
}

func (m metricsMiddlewareServer) AddPurchaseStatus(ctx context.Context, in PurchaseStatus) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddPurchaseStatus"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.AddPurchaseStatus(ctx, in)
	return
}

func (m metricsMiddlewareServer) AddTexture(ctx context.Context, in Texture) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddTexture"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.AddTexture(ctx, in)
	return
}

func (m metricsMiddlewareServer) AddProcess(ctx context.Context, in Process) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddProcess"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.AddProcess(ctx, in)
	return
}

func (m metricsMiddlewareServer) AddCraft(ctx context.Context, in Craft) (out NewParameterAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "AddCraft"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.AddCraft(ctx, in)
	return
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

func (m metricsMiddlewareServer) SearchWeightWithMaterialCode(ctx context.Context, in MaterialCode) (out WeightMaterialCodeAck, err error) {
	defer func(start time.Time) {
		method := []string{"method", "SearchWeightWithMaterialCode"}
		m.counter.With(method...).Add(1)
		m.histogram.With(method...).Observe(time.Since(start).Seconds())
	}(time.Now())
	out, err = m.next.SearchWeightWithMaterialCode(ctx, in)
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
