package stats

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/influxdata/influxdb/client/v2"
)

type influx struct {
	Config   client.HTTPConfig
	Database string

	Points []*client.Point
	Tags   []map[string]string
}

type Point struct {
	Fields map[string]interface{}
}

func (w *influx) WriteBwTest(r BpsResult) {
	f := map[string]interface{}{
		"BPS": r.Bps,
	}
	t := map[string]string{
		"testtag": "spaghett",
	}
	point, err := client.NewPoint(
		"bps",
		t,
		f,
	)
	errors.CheckError(err)
	w.Points = append(w.Points, point)
}

func (w *influx) WriteBwSummary(r BpsSummaryResult) {
	c, err := client.NewHTTPClient(w.Config)
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

func (w *influx) WriteReliabilityTest(r ReliabilityResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
func (w *influx) WriteReliabilitySummary(r ReliabilitySummaryResult) {
	fmt.Printf("Loss: %v, effective loss: %v\n", r.Loss, r.EffectiveLoss)
}
