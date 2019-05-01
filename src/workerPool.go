package WorkerPool

import time
import context
import golang.org/x/sync/semaphore

// define generic job type
type JobType func(workload interface{}) (interface{})


type interface WorkerPool {
  numWorkers int
  jobs chan<- JobType
  result <-chan interface{}
  AddJob(interface{})
  GetResult() interface{}
}

// ------- gorountine implementation ----------
type WorkerPoolThread struct {
  numWorkers int
  jobs chan<- JobType
  results <-chan interface{}
}

// new worker pool creates a new pool of workers where each worker will process using the provided func
func NewWorkerPoolThread(int num) &WorkerPoolThread {
  wp := &WorkerPool{
    num,
    make(chan<- interface{}),
    make(<-chan interface{})
  }
  for w:=1; w <= num; w++ {
    go func(){
      for j := range wp.jobs{
        results<-j()
      }
    }
  }
  return wp
}

func (WorkerPoolThread* wp) AddJob(JobType job) {
  wp.jobs <- job
}

func (WorkerPoolThread* wp) GetResult() interface{} {
  return <-wp.results
}




// -------- semaphore implementation -----------
type interface WorkerPoolSemaphore {
  numberWorkers int
  sem semaphore.NewWeighted
  jobs chan<- JobType
  results <-chan interface{}
  ctx Context
}

func NewWorkerPoolSemaphore(num int) &NewWorkerPoolSemaphore {
  return &WorkerPoolSemaphore{
    num,
    semaphore.NewWeighted(int64(num)),
    make(chan<- JobType),
    make(<-chan interface{}),
    context.TODO()
  }
}

// initialize the workpool
func (WorkerPoolSemaphore* wp) Init() {
  // keep executing jobs while there is job
  for j := range wp.jobs {
    // err when all resources are in use
    if err := sem.Acquire(ctx, 1); err != nil {
        log.Printf("Failed to acquire semaphore: %v", err)
        break
    }

    go func() {
      defer sem.Release(1)
      wp.results <- j()
    }
  }
}

func (WorkerPoolSemaphore* wp) AddJob(JobType job) {
  wp.jobs <- job
}

func (WorkerPoolSemaphore* wp) GetResult() interface{} {
  return <- wp.results
}
