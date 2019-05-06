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
	w.WritePoints()
}

func (w *Influx) WriteReliabilityTest(r ReliabilityResult) {
	f := map[string]interface{}{
		"Loss":          r.Loss,
		"EffectiveLoss": r.EffectiveLoss,
	}
	t := w.Tags
	point, err := client.NewPoint(
		"reliability",
		t,
		f,
	)
	errors.CheckError(err)
	w.Points = append(w.Points, point)
}
func (w *Influx) WriteReliabilitySummary(r ReliabilitySummaryResult) {
	w.WritePoints()
}

func (w *Influx) WritePoints() {
	fmt.Printf("Writing to influx.\n")
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
