package BoundedBuffer

import sync
import golang.org/x/sync/semaphore

// This implementation support concurrently read and write
type BoundedBuffer struct {
  r_idx_lock sync.Mutex
  w_idx_lock sync.Mutex
  full semaphore.NewWeighted
  empty semaphore.NewWeighted
  buffer []byte
  r_idx int
  w_idx int
/*
  empty_count semaphore.NewWeighted
  fill_count semaphore.NewWeighted
  buffer_lock sync.Mutex
  buffer Buffer
  */
}

func NewBuffer(size int) {
  return &BoundedBuffer{
    sync.Mutex{},
    sync.Mutex{},
    semaphore.NewWeighted(uint64(0)),
    semaphore.NewWeighted(uint64(size)),
    [size]byte,
    0,
    0,
  }
}

func (BoundedBuffer* bf) produce(item int) {
  // wait write until there is empty slot
  bf.full.Acquire()

  bf.w_idx_lock.Lock()
  bf.w_idx += 1
  bf.w_idx_lock.Unlock()

  bf.buffer[bf.w_idx] = item

  bf.empty.Release()

}

func (BoundedBuffer* bf) consume() int {
  // wait until there is an empty one
  bf.empty().Acquire()

  bf.r_idx_lock.Lock()
  bf.r_idx +=1
  bf.r_idx_lock.Unlock()

  ret_item = bf.buffer[bf.r_idx]

  bf.full.Release()
  return ret_item

}
