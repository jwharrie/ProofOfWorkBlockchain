package work_queue

// Each Worker instance executes Run() to start mining.
type Worker interface {
	Run() interface{}
}
/*
Struct containing two channels:
	Jobs: channel for Worker instances enter to mine.
	Results: each Worker instance sends results to this channel
*/
type WorkQueue struct {
	Jobs    chan Worker
	Results chan interface{}
}

// Create a new work queue capable of doing nWorkers simultaneous tasks, expecting to queue maxJobs tasks.
func Create(nWorkers uint, maxJobs uint) *WorkQueue {
	q := new(WorkQueue)
	q.Jobs = make(chan Worker, maxJobs)		// make new Worker channel that can contain up to maxJobs Workers
	q.Results = make(chan interface{})
	for i := uint(0); i < nWorkers; i++ {	// start goroutines
		go q.worker()
	}
	return q
}

// A worker goroutine that processes tasks from .Jobs unless .StopRequests has a message saying to halt now.
func (queue WorkQueue) worker() {
	for t := range queue.Jobs {
		queue.Results <- t.Run()
	}
	return
}

func (queue WorkQueue) Enqueue(work Worker) {
	queue.Jobs <- work
	return
}

func (queue WorkQueue) Shutdown() {
	close(queue.Jobs)
	for t := range queue.Jobs {
		_ = t
	}
	return
}
