//package BitSet
// thread safe 
// This is a dynanmic bitset that can expand

package bitset

import "fmt"

// This is a dynanmic bitset with arbitrary large range number
// each block stores 64 bit intger
type BitSet struct {
	blocks []uint64
}

func NewBitSet(size int) (*BitSet) {
	return &BitSet{
		blocks: make([]uint64, size),
	}
}

func (bs *BitSet) Check(position int) (bool) {
	// figure out the block chunk first
	block_num := position/64
	bs.checkBlocks(block_num)
	
	return ((1 << uint(position - block_num * 64)) & bs.blocks[block_num]) > 0
}

func (bs *BitSet) Set(position int) {
	// figure out the block chunk first
	block_num := position/64
	bs.checkBlocks(block_num)


	bs.blocks[block_num] = (1 << uint(position - block_num * 64)) | bs.blocks[block_num]
}

func (bs *BitSet) checkBlocks(blockNum int){
	for (blockNum >= len(bs.blocks)){
		bs.blocks = append(bs.blocks, 0)
	}
}

/*
func main(){
	bs := NewBitSet(0)
	bs.Set(100)
	bs.Set(0)
	fmt.Println(bs.Check(100))
	fmt.Println(bs.Check(200))
	bs.Set(0)
	fmt.Println(bs.Check(0))
}
*/
