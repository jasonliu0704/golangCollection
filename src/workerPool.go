package WorkerPool

type WorkerPool struct {
  numWorkers int
  jobs <-chan func(workload interface{}) (chan<- interface{})
}

// new worker pool creates a new pool of workers where each worker will process using the provided func
func NewWorkerPool(int num) &WorkerPool {
  wp := &WorkerPool{num}
  for w:=1; w <= num; w++ {
    go func(){
      for j := range wp.jobs{
        j()
      }
    }
  }
  return wp
}
