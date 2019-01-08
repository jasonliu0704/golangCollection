//package BitSet
// thread safe 

package bitset

import "fmt"

// This is a dynanmic bitset with arbitrary large range number
// each block stores 64 bit intger
type BitSet struct {
	blocks []uint64
}

func NewBitSet() (*BitSet) {
	return new(BitSet)
}

func (bs *BitSet) check(position int) (bool) {
	// figure out the block chunk first
	block_num := position/64
	bs.checkBlocks(block_num)
	
	return ((1 << uint(position - block_num * 64)) & bs.blocks[block_num]) > 0
}

func (bs *BitSet) set(position int) {
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
	bs := NewBitSet()
	bs.set(100)
	bs.set(0)
	fmt.Println(bs.check(100))
	fmt.Println(bs.check(200))
	bs.set(0)
	fmt.Println(bs.check(0))
}
*/
