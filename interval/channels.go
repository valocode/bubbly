package interval

import (
	"fmt"
	"time"

	"github.com/verifa/bubbly/api/core"
)

type Channels map[string]chan PipelineAction

type ChannelEvent struct {
	Interval time.Duration
	Cancel   string // custom string type
}

// UpdateResourceWorker will update a running interval Worker with an updated ResourceBlock
// for example, if a ResourceBlock is updated via POST, this will send the updated information
// to the ResourceWorker responsible.
func (w *ResourceWorker) UpdateResourceWorker(resBlock *core.ResourceBlock) error {
	if _, ok := w.WorkerChannels[resBlock.String()]; ok {
		w.WorkerChannels[resBlock.String()] <- PipelineAction{
			Action:        UpdatePipeline,
			ResourceBlock: resBlock,
		}
	} else {
		return fmt.Errorf("pipeline_run has not been initialized properly")
	}
	return nil
}

// DeleteChannel will delete and close the specified channel
func (pr *PipelineRun) DeleteChannel() {
	pr.Channel <- PipelineAction{
		Action:        StopPipeline,
		ResourceBlock: nil,
	}
}

// CloseChannels will close all active channels when all pipelines need to be closed
func (w *ResourceWorker) CloseChannels() {
	w.Context.Cancel()
	for channel := range w.WorkerChannels {
		delete(w.WorkerChannels, channel)
	}
}
