// Copyright 2023 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"gvisor.dev/gvisor/pkg/atomicbitops"
	"gvisor.dev/gvisor/pkg/log"
)

const snapshotBufferSize = 1000

var (
	// ProfilingMetricWriter is the output destination to which
	// ProfilingMetrics will be written to in CSV format.
	ProfilingMetricWriter *os.File
	// profilingMetricsStarted indicates whether StartProfilingMetrics has
	// been called.
	profilingMetricsStarted atomicbitops.Bool
	// stopProfilingMetrics is used to signal to the profiling metrics
	// goroutine to stop recording and writing metrics.
	stopProfilingMetrics chan bool = make(chan bool, 1)
	// defaultProfilingMetrics is the set of metrics used if no
	// profiling-metrics have been provided in the runsc flags.
	// These metrics are guaranteed to be registered and valid
	// (see condmetric_profiling.go).
	defaultProfilingMetrics []string
)

// StartProfilingMetrics checks the ProfilingMetrics runsc flags and creates
// goroutines responsible for outputting the profiling metric data.
//
// Precondition:
//   - All metrics are registered. Initialize/Disable has been called.
func StartProfilingMetrics(profilingMetrics string, profilingRate time.Duration) error {
	if !initialized.Load() {
		// Wait for initialization to complete to make sure that all
		// metrics are registered.
		return errors.New("metric initialization is not complete")
	}
	if ProfilingMetricWriter == nil {
		return errors.New("tried to initialize profiling metrics without log file")
	}
	if !profilingMetricsStarted.CompareAndSwap(0, 1) {
		return errors.New("profiling metrics have already been started")
	}

	var values []func(fieldValues ...*FieldValue) uint64
	header := strings.Builder{}
	header.WriteString("Time")
	numMetrics := 0
	recordMetric := func(name string, m customUint64Metric) {
		if len(m.fields) > 0 {
			// TODO(b/240280155): Add support for field values.
			log.Warningf("Will not profile metric '%s' because it has metric fields which are not supported")
			return
		}
		header.WriteRune(',')
		header.WriteString(name)
		values = append(values, m.value)
		numMetrics++
	}

	if len(profilingMetrics) > 0 {
		metrics := strings.Split(profilingMetrics, ",")

		for _, name := range metrics {
			name := strings.TrimSpace(name)
			m, ok := allMetrics.uint64Metrics[name]
			if !ok {
				return fmt.Errorf("given profiling metric name '%s' does not correspond to a registered Uint64 metric", name)
			}
			recordMetric(name, m)
		}
	} else {
		for _, name := range defaultProfilingMetrics {
			m, _ := allMetrics.uint64Metrics[name]
			recordMetric(name, m)
		}
		// Output equivalent flag in case user needs to narrow it down.
		log.Infof("A value for --profiling-metrics was not specified. Using '--profiling-metrics=%s'", strings.Join(defaultProfilingMetrics, ","))
	}

	if numMetrics == 0 {
		log.Warningf("No Profiling Metrics have been specified via -profiling-metrics or loaded at initialization time, even though a profiling-metrics-log file has been specified. If you forgot to compile the conditionally compiled metrics, use '--go_tag=condmetric_profiling' when compiling runsc.")
		return nil
	}

	header.WriteRune('\n')
	header.WriteRune('0')
	for i := 0; i < numMetrics; i++ {
		header.WriteString(",0")
	}
	header.WriteRune('\n')

	writeCh := make(chan profilingSnapshot)
	go collectProfilingMetrics(numMetrics, values, profilingRate, writeCh)
	go writeProfilingMetrics(header.String(), numMetrics, writeCh)

	return nil
}

type profilingSnapshot struct {
	// data is made up of lines like {timestamp,metric1,metric2,...}.
	data         []uint64
	numSnapshots int
}

// collectProfilingMetrics will send metrics to the writeCh until it receives a
// signal via the stopProfilingMetrics channel.
func collectProfilingMetrics(numMetrics int, values []func(fieldValues ...*FieldValue) uint64, profilingRate time.Duration, writeCh chan<- profilingSnapshot) {
	numEntries := numMetrics + 1 // to account for the timestamp
	snapshots := make([]uint64, numEntries*snapshotBufferSize)
	curSnapshot := 0
	startTime := time.Now()

collect:
	for {
		time.Sleep(profilingRate)
		timestamp := time.Now().Sub(startTime).Microseconds()

		base := curSnapshot * numEntries
		snapshots[base] = uint64(timestamp)
		for i := 1; i < numEntries; i++ {
			snapshots[base+i] = values[i-1]()
		}
		curSnapshot++

		select {
		case <-stopProfilingMetrics:
			writeCh <- profilingSnapshot{data: snapshots, numSnapshots: curSnapshot}
			break collect
		default:
		}

		if curSnapshot == snapshotBufferSize {
			writeCh <- profilingSnapshot{data: snapshots, numSnapshots: curSnapshot}
			curSnapshot = 0
			snapshots = make([]uint64, numEntries*snapshotBufferSize)
		}
	}

	close(writeCh)
}

func writeProfilingMetrics(header string, numMetrics int, snapshots <-chan profilingSnapshot) {
	io.WriteString(ProfilingMetricWriter, header)

	numEntries := numMetrics + 1
	for {
		s, ok := <-snapshots
		if !ok {
			break
		}

		out := strings.Builder{}
		for i := 0; i < s.numSnapshots; i++ {
			base := i * numEntries
			// Write the time
			out.WriteString(strconv.FormatUint(s.data[base], 10))
			// Then everything else
			for j := 1; j < numEntries; j++ {
				out.WriteRune(',')
				out.WriteString(strconv.FormatUint(s.data[base+j], 10))
			}
			out.WriteRune('\n')
		}

		io.WriteString(ProfilingMetricWriter, out.String())
	}
	ProfilingMetricWriter.Close()
}

// StopProfilingMetrics stops the profiling metrics goroutines. Call to make sure
// all metric data has been flushed.
func StopProfilingMetrics() {
	select {
	case stopProfilingMetrics <- true:
	default: // Stop signal was already sent
	}
}
