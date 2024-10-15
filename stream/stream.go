package stream

import (
	"context"
	"sync"

	"github.com/linkdata/certstreamui/certificate/v1"
)

func Stream(
	ctx context.Context,
	logOps []Operator,
	logSts []LogStatus,
	startIndex, batchSize, nWorkers int,
) (certCh <-chan *certificate.Batch, err error) {
	// Initialize the operators
	var logOperators []*LogOperator
	if logOperators, err = GetOperatorsFromArg(logOps); err == nil {

		// Create a channel for communication between log operators and the sink
		fromLogsToSink := make(chan *certificate.Batch, nWorkers)
		certCh = fromLogsToSink

		// Start streaming logs from each operator
		go func() {
			var wg sync.WaitGroup
			defer close(fromLogsToSink)
			for _, logOperator := range logOperators {
				// Initialize streams for each operator
				logOperator.InitStreams(logSts, batchSize, nWorkers, startIndex)
				// Run this operator's streams in a goroutine
				wg.Add(1)
				go func(logOp *LogOperator) {
					defer wg.Done()
					logOp.RunStreams(ctx, fromLogsToSink)
				}(logOperator)
			}
			// Wait for operators to finish
			wg.Wait()
		}()
	}
	return
}
