/*
The read write Lock allow one write and multiple reader
read and write are exclusive
*/

// enhancement, to prevent both read and write starvation,
// use another mutex as a service queue
package RWLock

import sync

type RWLock struct {
  // number of readers
  readCount int
  // mutual exclusion to readcount
  read_count_mutex sync.Mutex
  // mutual exclusion to read/write
  rw_mutex sync.Mutex
  // prevent rw starvation
  serviceQueue syc.Mutex
}

func (RWLock* l) ReadLock() {

  serviceQueue.Lock()
  // lock readcount so we can check readcount == 1 to sync with writer
  l.read_count_mutex.Lock()
  serviceQueue.Unlock()
  // sync up with writer when we first start to read
  if(l.readCount == 1) {
    rw_mutex.Lock()
  }
  l.read_count_mutex.Unlock()
}

func (RWLock* l) ReadUnlock() {
  l.read_count_mutex.Lock()
  l.readCount -= 1
  // can only start writing when no one is reading
  if(l.readCount == 0) {
    rw_mutex.Unlock()
  }
  l.read_count_mutex.Unlock()

}

func (RWLock* l) WriteLock() {
  serviceQueue.Lock()
  rw_mutex.Lock()
  serviceQueue.Unlock()
}

func (RWLock* l) WriteUnlock() {
  rw_mutex.Unlock()
}
