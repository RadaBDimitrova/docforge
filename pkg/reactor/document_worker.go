package reactor

import (
	"context"
	"fmt"

	"github.com/gardener/docforge/pkg/api"
	"github.com/gardener/docforge/pkg/jobs"
	"github.com/gardener/docforge/pkg/processors"
	"github.com/gardener/docforge/pkg/resourcehandlers"
	utilnode "github.com/gardener/docforge/pkg/util/node"
	"github.com/gardener/docforge/pkg/writers"
)

// Reader reads the bytes data from a given source URI
type Reader interface {
	Read(ctx context.Context, source string) ([]byte, error)
}

// DocumentWorker implements jobs#Worker
type DocumentWorker struct {
	writers.Writer
	Reader
	processors.Processor
	NodeContentProcessor *NodeContentProcessor
	localityDomain       localityDomain
}

// DocumentWorkTask implements jobs#Task
type DocumentWorkTask struct {
	Node *api.Node
	// localityDomain localityDomain
}

// GenericReader is generic implementation for Reader interface
type GenericReader struct {
	ResourceHandlers resourcehandlers.Registry
}

// Read reads from the resource at the source URL delegating the
// the actual operation to a suitable resource hadler
func (g *GenericReader) Read(ctx context.Context, source string) ([]byte, error) {
	if handler := g.ResourceHandlers.Get(source); handler != nil {
		return handler.Read(ctx, source)
	}
	return nil, fmt.Errorf("failed to get handler")
}

// Work implements Worker#Work function
func (w *DocumentWorker) Work(ctx context.Context, task interface{}, wq jobs.WorkQueue) *jobs.WorkerError {
	if task, ok := task.(*DocumentWorkTask); ok {

		var sourceBlob []byte

		if len(task.Node.ContentSelectors) > 0 {
			// Write the document content
			blobs := make(map[string][]byte)
			for _, content := range task.Node.ContentSelectors {
				sourceBlob, err := w.Reader.Read(ctx, content.Source)
				if len(sourceBlob) == 0 {
					fmt.Printf("No bytes read from source %s\n", content.Source)
					continue
				}
				if err != nil {
					return jobs.NewWorkerError(err, 0)
				}
				newBlob, err := w.NodeContentProcessor.ReconcileLinks(ctx, task.Node, content.Source, sourceBlob)
				if err != nil {
					return jobs.NewWorkerError(err, 0)
				}
				blobs[content.Source] = newBlob
			}

			if len(blobs) == 0 {
				return nil
			}

			for _, blob := range blobs {
				sourceBlob = append(sourceBlob, blob...)
			}

			var err error
			if w.Processor != nil {
				if sourceBlob, err = w.Processor.Process(sourceBlob, task.Node); err != nil {
					return jobs.NewWorkerError(err, 0)
				}
			}
		}

		path := utilnode.Path(task.Node, "/")
		if err := w.Writer.Write(task.Node.Name, path, sourceBlob, task.Node); err != nil {
			return jobs.NewWorkerError(err, 0)
		}
	}
	return nil
}
