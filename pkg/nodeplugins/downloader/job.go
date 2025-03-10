// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package downloader

import (
	"context"
	"fmt"
	"sync"

	"github.com/gardener/docforge/pkg/registry"
	"github.com/gardener/docforge/pkg/workers/taskqueue"
	"github.com/gardener/docforge/pkg/writers"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate -header ../../../license_prefix.txt

// Interface encapsulates activities for asynchronous and parallel scheduling and download of resources
//
//counterfeiter:generate . Interface
type Interface interface {
	// Schedule is a typesafe wrapper for enqueuing download tasks. An error is returned if scheduling fails.
	Schedule(source string, destinationPath string) error
}

type downloadScheduler struct {
	*ResourceDownloadWorker
	queue taskqueue.Interface
}

// New create a DownloadScheduler to schedule download resources
func New(workerCount int, failFast bool, wg *sync.WaitGroup, registry registry.Interface, writer writers.Writer) (Interface, taskqueue.QueueController, error) {
	dWorker, err := NewDownloader(registry, writer)
	if err != nil {
		return nil, nil, err
	}
	queue, err := taskqueue.New("Download", workerCount, dWorker.execute, failFast, wg)
	if err != nil {
		return nil, nil, err
	}
	downloader := &downloadScheduler{
		dWorker,
		queue,
	}
	return downloader, queue, nil
}

// Schedule enqueues and resource link for download
func (ds *downloadScheduler) Schedule(source string, destinationPath string) error {
	task := &downloadTask{source, destinationPath}
	if !ds.queue.AddTask(task) {
		return fmt.Errorf("scheduling download of %s failed", task.source)
	}
	return nil
}

func (d *ResourceDownloadWorker) execute(ctx context.Context, task interface{}) error {
	dt, ok := task.(*downloadTask)
	if !ok {
		return fmt.Errorf("incorrect download task: %T", task)
	}
	return d.Download(ctx, dt.source, dt.destinationPath)
}

// DownloadTask holds information for source and destinationPath of linked document resources
type downloadTask struct {
	source          string
	destinationPath string
}
