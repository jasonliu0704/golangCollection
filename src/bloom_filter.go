package main

import "fmt"
import "math"
import "hash"
import "murmur3"
import "bitset"

type BloomFilter struct {
	bitset *BitSet // bitset to store bit
	k int // number of hash functions
	m int // size of the bitset
	n int // expected number of elements to be stores
	fp float32 // false positive probability
	hashFuncs[]hash.Hash64 // k hash functions
}

// m = - (nlnP)/(ln2)^2
func calcM(p float32){
	return - (n * math.Log(p))/math.Pow(math.Log(2), 2)
}

func calcK(m, n int){
	return (m/n) * math.Log(2) 
}

func NewBloomFilter(fp float32, n int){
	// default to 3 hash functions
	m := calcM(fp)
	k := calcK(m, n)
	var hashFuncs [k]hash.Hash64
	for i := 0; i < k; i++ {
    	hashFuncs[i] = murmur3.New64WithSeed(i)
	}
	return new(BloomFilter{
		bitset: bitset.NewBitSet(),
		k: k ,
		m: calcM(fp),
		n: n,
		fp: fp,
		hashfuncs: hashFuncs,
	}	
}

func (bf *BloomFilter) Set (item int){
	// calculate k hash first, and set each hash position 
	for _, h :=  range bf.hashFuncs {
		h.Write([]byte(item))
		bf.bitset.Set(h.Sum() % m)
		h.Reset()
	}
}

func (bf BloomFilter) Check (item int) (bool){
	for _, h := range bf.hashFuncs {
		h.Write([]byte(item))
		if(!bf.bitset.Check(h.Sum() % m)){
			h.Reset()
			return false
		}
		h.Reset()
		return true
	}
}

func main() {
	bf := NewBloomFilter(0.2, 100)
	bf.Set(100)
	bf.Set(10000000)
	bf.Set(200)
	bf.Set(0)
	bf.Set(99)
	fmt.Println(bf.Check(1))
	fmt.Pritnln(bf.Check(2))
	fmt.Println(bf.Check(11312314211312))
	fmt.Println(bf.Check(7))
	fmt.Println(bf.Check(77))
	fmt.Println(bf.Check(100))
	fmt.Println(bf.Check(99))
}

