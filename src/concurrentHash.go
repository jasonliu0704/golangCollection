I always visualize write locks as a funnel where concurrent execution is forced into a serial mode. You can build an awesome system capable of sustaining tens of thousands of concurrent processes, yet have a badly thought out lock that brings the entire system to a grind. Think of a highway reduced to a single lane.

If you want concurrent access to a hash table, there's one simple trick you can use to minimize the impact of a write lock: shard it. Before we look at the solution, let's look at the problem. Imagine you're building a simple cache backed by a map:

import (
  "sync"
)

type Cache struct {
  items map[string][]byte
  lock *sync.RWLock
}

func New() *Cache {
  return &Cache{
    items: make(map[string][]byte, 16384),
    lock: new(sync.RWLock),
  }
}

func (c *Cache) Get(key string) []byte {
  c.lock.RLock()
  defer c.lock.RUnlock()
  return c.items[key]
}

func (c *Cache) Set(key string, data []byte) {
  c.lock.Lock()
  defer c.lock.Unlock()
  c.items[key] = data
}
The problem with the above code is that any call to Set is not only serialized, but also blocks calls to Get. In this case, that might seem trivial. With real code though, you might need to do more than set a value, like checking for and freeing up sufficient space.

The solution to this problem is to shard our hash table. This is implemented by using a hashtable of hashtables, where each inner hash table has its own lock and represents a more granular space:

import (
  "fmt"
  "sync"
  "crypto/sha1"
)

type Cache struct {
  shards map[string]*CacheShard
  lock *sync.RWMutex
}

type CacheShard struct {
  items map[string][]byte
  lock *sync.RWMutex
}

func New() *Cache {
  return &Cache{
    shards: make(map[string]*CacheShard, 256),
    lock: new(sync.RWMutex),
  }
}

func (c *Cache) Get(key string) []byte {
  shard, ok := c.GetShard(key, false)
  if ok == false { return nil }
  shard.lock.RLock()
  defer shard.lock.RUnlock()
  return shard.items[key]
}

func (c *Cache) Set(key string, data []byte) {
  shard, _ := c.GetShard(key, true)
  shard.lock.Lock()
  defer shard.lock.Unlock()
  shard.items[key] = data
}

func (c *Cache) GetShard(key string, create bool) (shard *CacheShard, ok bool) {
  hasher := sha1.New()
  hasher.Write([]byte(key))
  shardKey :=  fmt.Sprintf("%x", hasher.Sum(nil))[0:2]

  c.lock.RLock()
  shard, ok = c.shards[shardKey]
  c.lock.RUnlock()

  if ok || !create { return }

  //only time we need to write lock
  c.lock.Lock()
  defer c.lock.Unlock()
  //check again in case the group was created in this short time
  shard, ok = c.shards[shardKey]
  if ok { return }

  shard = &CacheShard {
    items: make(map[string][]byte, 2048),
    lock: new(sync.RWMutex),
  }
  c.shards[shardKey] = shard
  ok = true
  return
}
By evenly distributing our keys across 256 shards (16^2) we effectively unclog our bottleneck. The only time we hold a global write lock is when the shard itself doesn't exist. Over time, this will be 0. As a matter of fact, since we know all possible shard keys ahead of time, we can pre-allocate shards and eliminate the need for the outer-most lock (since we'll only ever read from it). This not only improve performance, it also makes the code much simpler:

import (
  "fmt"
  "sync"
  "crypto/sha1"
)

type Cache map[string]*CacheShard

type CacheShard struct {
  items map[string][]byte
  lock *sync.RWMutex
}

func New() Cache {
  c := make(Cache, 256)
  for i := 0; i < 256; i++ {
    c[fmt.Sprintf("%02x", i)] = &CacheShard{
      items: make(map[string][]byte, 2048),
      lock: new(sync.RWMutex),
    }
  }
  return c
}

func (c Cache) Get(key string) []byte {
  shard := c.GetShard(key)
  shard.lock.RLock()
  defer shard.lock.RUnlock()
  return shard.items[key]
}

func (c Cache) Set(key string, data []byte) {
  shard := c.GetShard(key)
  shard.lock.Lock()
  defer shard.lock.Unlock()
  shard.items[key] = data
}

func (c Cache) GetShard(key string) (shard *CacheShard) {
  hasher := sha1.New()
  hasher.Write([]byte(key))
  shardKey :=  fmt.Sprintf("%x", hasher.Sum(nil))[0:2]
  return c[shardKey]
}
It's worth pointing out that some libraries, like Java's ConcurrentHashMap, does this internally already.

You'll need to decide if the extra complexity, the extra memory overhead and the SHA1 calculation is worth it. What you're doing while locked as well as how often you'll be locking, will be major factors. Though, to be honest, I think in most cases, the concurrency gains will outweigh any raw throughput concerns.
