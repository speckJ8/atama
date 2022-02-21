package device

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
	blockSize uint
	setSize   uint
	setCount  uint
	// The length of the the data array
	// should equal blockSize*setSize*setCount
	blocks []CacheBlock
	typ    CacheType
}

type CacheBlock struct {
	Used    bool
	Active  bool
	Address uint
	Data    []byte
}

func NewCache(blockSize, setCount, setSize uint) Cache {
	if blockSize%8 != 0 {
		blockSize = blockSize - (blockSize % 8) + 8
	}
	sets := make([]CacheBlock, setCount*setSize)
	for s := range sets {
		sets[s] = CacheBlock{
			Active: false,
			Data:   make([]byte, blockSize),
		}
	}
	return Cache{
		blockSize: blockSize,
		setCount:  setCount,
		setSize:   setSize,
		blocks:    sets,
	}
}

// TODO: write back the cache line that was ejected
func (c *Cache) Populate(address uint, mem *Memory) MemoryAccessStatus {
	address = address - address%c.blockSize
	statusChannel, dataChannel := mem.Read(address, c.blockSize)
	status := <-statusChannel
	if status == MemoryAccessReadDone {
		block := <-dataChannel
		set := address % c.setCount
		s := set
		for ; s < set+c.setSize; s++ {
			if !c.blocks[s].Active {
				c.blocks[s].Active = true
				c.blocks[s].Used = true
				c.blocks[s].Data = block
				c.blocks[s].Address = address
				break
			} else if !c.blocks[s].Used {
				c.blocks[s].Used = true
				c.blocks[s].Data = block
				c.blocks[s].Address = address
				break
			}
		}
		if s == set+c.setSize {
			s = set
			for ; s < set+c.setSize; s++ {
				c.blocks[s].Used = false
			}
			c.blocks[set].Used = true
			c.blocks[set].Data = block
			c.blocks[set].Address = address
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
	return CacheAccessWriteDone
}

func (c *Cache) SetByte(address uint, data Byte) CacheAccessStatus {
	block := c.getBlock(address)
	if block == nil {
		return CacheAccessMiss
	}
	start := c.globalAddressToBlockAddres(address)
	block.Data[start] = data
	return CacheAccessWriteDone
}

func (c *Cache) getBlock(address uint) *CacheBlock {
	address = address - address%c.blockSize
	set := address % c.setCount
	for s := set; s < set+c.setSize; s++ {
		if c.blocks[s].Active && c.blocks[s].Address == address {
			return &c.blocks[s]
		}
	}
	return nil
}

func (c *Cache) globalAddressToBlockAddres(address uint) int {
	return int(address % c.blockSize)
}
