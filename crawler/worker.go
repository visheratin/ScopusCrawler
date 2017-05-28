package crawler

import "fmt"

type Worker struct {
	Work        chan SearchRequest
	WorkerQueue chan chan SearchRequest
}

func (worker *Worker) Start() {
	go func() {
		for {
			worker.WorkerQueue <- worker.Work
			select {
			case work := <-worker.Work:
				fmt.Println(work.SourceName)
				for k, v := range work.Fields {
					fmt.Println(k, v)
				}
			}
		}
	}()
}
