package stats

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/influxdata/influxdb/client/v2"
)

type Influx struct {
	HTTPConfig client.HTTPConfig
	Database   string

	Points []*client.Point
	Tags   map[string]string
}

type Config struct {
}

type Point struct {
	Fields map[string]interface{}
}

func (w *Influx) WriteBwTest(r BpsResult) {
	f := map[string]interface{}{
		"BPS": r.Bps,
	}
	t := w.Tags
	point, err := client.NewPoint(
		"bps",
		t,
		f,
	)
	errors.CheckError(err)
	w.Points = append(w.Points, point)
}

func (w *Influx) WriteBwSummary(r BpsSummaryResult) {
	c, err := client.NewHTTPClient(w.HTTPConfig)
	errors.CheckError(err)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  w.Database,
		Precision: "s",
	})

	for _, p := range w.Points {
		bp.AddPoint(p)
	}
	err = c.Write(bp)
	errors.CheckError(err)
}

func (w *Influx) WriteReliabilityTest(r ReliabilityResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
func (w *Influx) WriteReliabilitySummary(r ReliabilitySummaryResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
