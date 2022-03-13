package device

import "math"

type CacheType = Byte
type CacheAccessStatus = Byte

const (
	CacheWriteBack    CacheType = iota
	CacheWriteThrough CacheType = iota
)

const (
	CacheAccessReadDone  CacheAccessStatus = iota
	CacheAccessWriteDone CacheAccessStatus = iota
	CacheAccessUnaligned CacheAccessStatus = iota
	CacheAccessMiss      CacheAccessStatus = iota
)

type Cache struct {
	BlockSize     uint
	blockSizeBits uint
	SetSize       uint
	SetCount      uint
	// The length of the the data array
	// should equal BlockSize*SetSize*SetCount
	Blocks []CacheBlock
}

type CacheBlock struct {
	// `Valid` indicates whether or not the block
	// actually contains data referring to a memory location
	// as opposed to just being uninitialized
	Valid bool
	// `RecentlyUsed` is set whenever a block is read or written and
	// is used in the eviction logic
	RecentlyUsed bool
	Dirty        bool
	Address      uint
	Data         []byte
}

func NewCache(BlockSize, SetCount, SetSize uint) Cache {
	sets := make([]CacheBlock, SetCount*SetSize)
	for s := range sets {
		sets[s] = CacheBlock{
			Data: make([]byte, BlockSize),
		}
	}
	return Cache{
		BlockSize:     BlockSize,
		blockSizeBits: uint(math.Log2(float64(BlockSize))),
		SetCount:      SetCount,
		SetSize:       SetSize,
		Blocks:        sets,
	}
}

func (c *Cache) Size() uint {
	return c.SetCount * c.SetSize * c.BlockSize
}

// TODO: write back the cache line that was ejected
func (c *Cache) Populate(address uint, mem *Memory) MemoryAccessStatus {
	address = address - address%c.BlockSize
	statusChannel, dataChannel := mem.Read(address, c.BlockSize)
	status := <-statusChannel
	if status == MemoryAccessReadDone {
		block := <-dataChannel
		set := (address >> c.blockSizeBits) % c.SetCount
		// look for a block in the set where to put
		// the contents fetched from memory
		i := uint(0)
		for ; i < c.SetSize; i++ {
			a := set + i*c.SetCount
			if !c.Blocks[a].RecentlyUsed {
				c.Blocks[a].RecentlyUsed = true
				c.Blocks[a].Valid = true
				c.Blocks[a].Data = block
				c.Blocks[a].Address = address
				break
			} else if !c.Blocks[a].Valid {
				c.Blocks[a].Valid = true
				c.Blocks[a].Data = block
				c.Blocks[a].Address = address
				break
			}
		}
		if i == c.SetSize {
			// there were no invalid and no stale blocks
			// in this set so we will have write to the
			// first block and make every other block stale
			i = uint(0)
			for ; i < c.SetSize; i++ {
				a := set + i*c.SetCount
				c.Blocks[a].RecentlyUsed = false
			}
			c.Blocks[set].Valid = true
			c.Blocks[set].Data = block
			c.Blocks[set].Address = address
		}
	}
	return status
}

func (c *Cache) GetQuadWord(address uint) (QWord, CacheAccessStatus) {
	if address%8 != 0 {
		return ZeroQWord, CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return ZeroQWord, CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	var qword QWord
	for i := 0; i < 8; i++ {
		qword[i] = block.Data[start+i]
	}
	return qword, CacheAccessReadDone
}

func (c *Cache) GetDoubleWord(address uint) (DWord, CacheAccessStatus) {
	if address%4 != 0 {
		return ZeroDWord, CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return ZeroDWord, CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	dword := DWord{
		block.Data[start],
		block.Data[start+1],
		block.Data[start+2],
		block.Data[start+3],
	}
	return dword, CacheAccessReadDone
}

func (c *Cache) GetWord(address uint) (Word, CacheAccessStatus) {
	if address%2 != 0 {
		return ZeroWord, CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return ZeroWord, CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	return Word{block.Data[start], block.Data[start+1]}, CacheAccessReadDone
}

func (c *Cache) GetByte(address uint) (Byte, CacheAccessStatus) {
	block := c.getBlock(address)
	if block == nil {
		return 0, CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	return Byte(block.Data[start]), CacheAccessReadDone
}

func (c *Cache) SetQuadWord(address uint, data QWord) CacheAccessStatus {
	if address%8 != 0 {
		return CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	for i := 0; i < 8; i++ {
		block.Data[start+i] = data[i]
	}
	block.Dirty = true
	return CacheAccessWriteDone
}

func (c *Cache) SetDoubleWord(address uint, data DWord) CacheAccessStatus {
	if address%4 != 0 {
		return CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	block.Data[start] = data[0]
	block.Data[start+1] = data[1]
	block.Data[start+2] = data[2]
	block.Data[start+3] = data[3]
	block.Dirty = true
	return CacheAccessWriteDone
}

func (c *Cache) SetWord(address uint, data Word) CacheAccessStatus {
	if address%2 != 0 {
		return CacheAccessUnaligned
	}
	block := c.getBlock(address)
	if block == nil {
		return CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	block.Data[start] = data[0]
	block.Data[start+1] = data[1]
	block.Dirty = true
	return CacheAccessWriteDone
}

func (c *Cache) SetByte(address uint, data Byte) CacheAccessStatus {
	block := c.getBlock(address)
	if block == nil {
		return CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	block.Data[start] = data
	block.Dirty = true
	return CacheAccessWriteDone
}

func (c *Cache) getBlock(address uint) *CacheBlock {
	set := (address >> c.blockSizeBits) % c.SetCount
	address = address - address%c.BlockSize
	for i := uint(0); i < c.SetSize; i++ {
		a := set + i*c.SetCount
		if c.Blocks[a].Valid && c.Blocks[a].Address == address {
			return &c.Blocks[a]
		}
	}
	return nil
}

func (c *Cache) globalAddressToBlockAddres(address uint) int {
	return int(address % c.BlockSize)
}
