package stats

import (
	"fmt"
	"github.com/adamb/netpupper/errors"
	"github.com/influxdata/influxdb/client/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Influx struct {
	Config Config `yaml:"Influx"`

	Points []*client.Point
	Tags   []map[string]string
}

type Config struct {
	HTTPConfig *client.HTTPConfig
	Database   string
}

type Point struct {
	Fields map[string]interface{}
}

func (w *Influx) Configure(cf string) bool {
	var serverFile string
	if len(cf) > 0 {
		serverFile = cf
	} else if len(os.Getenv("NETP_CONFIG")) > 0 {
		serverFile = os.Getenv("NETP_CONFIG")
	}
	// Do you have to do this, for real?
	w.Config = Config{
		HTTPConfig: &client.HTTPConfig{},
	}
	// If the yaml file exists
	if _, err := os.Stat(serverFile); err == nil {
		data, err := ioutil.ReadFile(serverFile)
		errors.CheckError(err)

		err = yaml.Unmarshal(data, w)
		errors.CheckError(err)
	}
	return false
}

func (w *Influx) WriteBwTest(r BpsResult) {
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

func (w *Influx) WriteBwSummary(r BpsSummaryResult) {
	c, err := client.NewHTTPClient(*w.Config.HTTPConfig)
	errors.CheckError(err)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  w.Config.Database,
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
