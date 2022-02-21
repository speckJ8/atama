package device

const WordSize = 2
const DWordSize = 2 * WordSize
const QWordSize = 4 * WordSize

type QWord = [QWordSize]byte
type DWord = [DWordSize]byte
type Word = [WordSize]byte
type Byte = byte

var ZeroQWord QWord
var ZeroDWord DWord
var ZeroWord Word
var ZeroByte Byte = 0
