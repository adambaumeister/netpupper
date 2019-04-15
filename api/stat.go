package api

import (
	"encoding/json"
	"github.com/adamb/netpupper/errors"
	"github.com/adamb/netpupper/perf/stats"
	"net/http"
)

type ApiCollector struct {
	rw http.ResponseWriter
}

func (a *ApiCollector) SetResponse(w http.ResponseWriter) {
	a.rw = w
}

func (a *ApiCollector) WriteBwTest(r stats.BpsResult) {
	// Do nothing - we don't send periodic updates to an API test collector.
}

func (a *ApiCollector) WriteBwSummary(r stats.BpsSummaryResult) {
	b, err := json.Marshal(r)
	errors.CheckError(err)
	a.rw.Write(b)
}

func (a *ApiCollector) WriteReliabilityTest(r stats.ReliabilityResult) {
	// Currently does nothing - we don't send periodic updates
}
func (a *ApiCollector) WriteReliabilitySummary(r stats.ReliabilitySummaryResult) {
	b, err := json.Marshal(r)
	errors.CheckError(err)
	a.rw.Write(b)
}
