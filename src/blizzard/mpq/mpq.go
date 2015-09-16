// Package mpq reads data from Blizzard MPQ files.
package mpq

import (
	"bufio"
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func read8(r io.Reader) (val uint8) {
	var buf [1]byte
	_, err := r.Read(buf[:])
	check(err)
	return buf[0]
}

func read16(r io.Reader) (val uint16) {
	check(binary.Read(r, binary.LittleEndian, &val))
	return
}

func read32(r io.Reader) (val uint32) {
	check(binary.Read(r, binary.LittleEndian, &val))
	return
}

func read64(r io.Reader) (val uint64) {
	check(binary.Read(r, binary.LittleEndian, &val))
	return
}

var cryptTable [0x500]uint32

func init() {
	seed := uint32(0x00100001)
	for index1 := 0; index1 < 0x100; index1++ {
		offset := index1
		for i := 0; i < 5; i++ {
			seed = (seed*125 + 3) % 0x2AAAAB
			temp1 := (seed & 0xFFFF) << 0x10
			seed = (seed*125 + 3) % 0x2AAAAB
			temp2 := (seed & 0xFFFF)

			cryptTable[offset] = temp1 | temp2
			offset += 0x100
		}
	}
}

// HashType identifies which variant of the hash function to use.
type HashType int

const (
	HashTableOffset = HashType(0)
	HashNameA       = HashType(1)
	HashNameB       = HashType(2)
	HashFileKey     = HashType(3)
)

func toupper(c uint8) uint8 {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}

// Hash is the hash function used within MPQ files.
func Hash(in string, hashType HashType) uint32 {
	seed1 := uint32(0x7fed7fed)
	seed2 := uint32(0xeeeeeeee)
	for i := range in {
		c := toupper(uint8(in[i]))
		seed1 = cryptTable[(int(hashType)*0x100)+int(c)] ^ (seed1 + seed2)
		seed2 = uint32(c) + seed1 + seed2 + (seed2 << 5) + 3
	}
	return seed1
}

type decrypter struct {
	r        io.Reader // underlying reader
	key      uint32
	seed     uint32
	extraBuf [4]byte
	extra    []byte
}

func newDecrypter(r io.Reader, key uint32) *decrypter {
	return &decrypter{
		r:    r,
		key:  key,
		seed: 0xeeeeeeee,
	}
}

func (d *decrypter) Read(buf []byte) (n int, err error) {
	n = 0
	for n < len(buf) {
		if d.extra != nil {
			copied := copy(buf, d.extra)
			n += copied
			d.extra = d.extra[copied:]
			if len(d.extra) == 0 {
				d.extra = nil
			}
			continue
		}

		var block uint32
		err = binary.Read(d.r, binary.LittleEndian, &block)
		if err != nil {
			return
		}

		d.seed += cryptTable[0x400+d.key&0xFF]
		block = block ^ (d.key + d.seed)
		d.key = ((^d.key << 0x15) + 0x11111111) | (d.key >> 0xB)
		d.seed = block + d.seed + (d.seed << 5) + 3
		if len(buf) >= 4 {
			binary.LittleEndian.PutUint32(buf[n:], block)
			n += 4
		} else {
			d.extra = d.extraBuf[:]
			binary.LittleEndian.PutUint32(d.extra, block)
		}
	}
	return
}

// Reader reads an MPQ file.
type Reader struct {
	*os.File
	userData   userData
	header     header
	hashTable  []hashEntry
	blockTable []blockEntry
}

type userData struct {
	size      uint32
	headerOfs uint32
	unk       uint32
}

func (r *Reader) readUserData() {
	u := &r.userData
	u.size = read32(r)
	u.headerOfs = read32(r)
	u.unk = read32(r)
}

type header struct {
	headerSize, archiveSize             uint32
	version, blockSize                  uint16
	hashTableOfs, blockTableOfs         uint32
	hashTableEntries, blockTableEntries uint32

	// v2 info
	extendedBlockTableOfs uint64
	hiHashTableOfs        uint16
	hiBlockTableOfs       uint16

	// v3 info
	archiveSize64            uint64
	betTablePos, hetTablePos uint64
}

func (r *Reader) readHeader() {
	h := &r.header
	h.headerSize = read32(r)
	h.archiveSize = read32(r)
	h.version = read16(r)
	h.blockSize = read16(r)

	h.hashTableOfs = read32(r)
	h.blockTableOfs = read32(r)
	h.hashTableEntries = read32(r)
	h.blockTableEntries = read32(r)

	if h.version >= 2 {
		h.extendedBlockTableOfs = read64(r)
		h.hiHashTableOfs = read16(r)
		h.hiBlockTableOfs = read16(r)
	}

	if h.version >= 3 {
		h.archiveSize64 = read64(r)
		h.betTablePos = read64(r)
		h.hetTablePos = read64(r)
	}

	if h.version >= 4 {
		panic("unimplemented version")
	}
}

type hashEntry struct {
	pathHashA  uint32
	pathHashB  uint32
	language   uint16
	platform   uint16
	blockIndex uint32
}

func (r *Reader) readHashTable() {
	ofs := r.userData.headerOfs + r.header.hashTableOfs
	_, err := r.Seek(int64(ofs), 0)
	check(err)

	r.hashTable = make([]hashEntry, r.header.hashTableEntries)
	d := newDecrypter(r, Hash("(hash table)", HashFileKey))
	for i := uint32(0); i < r.header.hashTableEntries; i++ {
		he := &r.hashTable[i]
		he.pathHashA = read32(d)
		he.pathHashB = read32(d)
		he.language = read16(d)
		he.platform = read16(d)
		he.blockIndex = read32(d)
	}
}

const (
	BlockFlagFile           uint32 = 1 << 31
	BlockFlagCheckSums      uint32 = 1 << 26
	BlockFlagDeletionMarker uint32 = 1 << 25
	BlockFlagSingleUnit     uint32 = 1 << 24
	BlockFlagCompressed     uint32 = 1 << 9
	BlockFlagImploded       uint32 = 1 << 8
)

type flagName struct {
	flag uint32
	name string
}

var flagNames []flagName = []flagName{
	{BlockFlagFile, "file"},
	{BlockFlagCheckSums, "checksums"},
	{BlockFlagDeletionMarker, "deletion marker"},
	{BlockFlagSingleUnit, "single unit"},
	{BlockFlagCompressed, "compressed"},
	{BlockFlagImploded, "imploded"},
}

type blockEntry struct {
	offset   uint32
	size     uint32
	fileSize uint32
	flags    uint32
}

func (be *blockEntry) flagsString() string {
	flags := be.flags
	flagList := []string{}
	for _, fn := range flagNames {
		if flags&fn.flag != 0 {
			flagList = append(flagList, fn.name)
			flags = flags & ^fn.flag
		}
	}

	if flags > 0 {
		flagList = append(flagList, fmt.Sprintf("[unknown flags %x]", flags))
	}
	return strings.Join(flagList, ", ")
}

func (r *Reader) readBlockTable() {
	ofs := r.userData.headerOfs + r.header.blockTableOfs
	_, err := r.Seek(int64(ofs), 0)
	check(err)

	r.blockTable = make([]blockEntry, r.header.blockTableEntries)
	d := newDecrypter(r, Hash("(block table)", HashFileKey))
	for i := uint32(0); i < r.header.blockTableEntries; i++ {
		he := &r.blockTable[i]
		he.offset = read32(d)
		he.size = read32(d)
		he.fileSize = read32(d)
		he.flags = read32(d)
	}
}

type het struct {
	version  uint32
	dataSize uint32

	tableSize      uint32
	maxFileCount   uint32
	hashTableSize  uint32
	hashEntrySize  uint32
	totalIndexSize uint32
	indexSizeExtra uint32
	indexSize      uint32
	blockTableSize uint32

	table []byte
}

func (r *Reader) readHET() *het {
	ofs := uint64(r.userData.headerOfs) + r.header.hetTablePos
	var buf [4]byte
	_, err := r.Seek(int64(ofs), 0)
	check(err)
	_, err = r.Read(buf[:])
	check(err)
	if string(buf[:]) != "HET\x1a" {
		panic(fmt.Errorf("bad HET signature %q", buf))
	}

	h := het{}
	h.version = read32(r)
	h.dataSize = read32(r)

	d := newDecrypter(r, Hash("(hash table)", HashFileKey))
	h.tableSize = read32(d)
	h.maxFileCount = read32(d)
	h.hashTableSize = read32(d)
	h.hashEntrySize = read32(d)
	h.totalIndexSize = read32(d)
	h.indexSizeExtra = read32(d)
	h.indexSize = read32(d)
	h.blockTableSize = read32(d)
	log.Printf("%#v", h)
	h.table = make([]byte, h.hashTableSize)
	_, err = d.Read(h.table)
	check(err)
	log.Printf("%#v", h)
	return &h
}

func (r *Reader) findFile(name string) *hashEntry {
	index := Hash(name, HashTableOffset) & (r.header.hashTableEntries - 1)
	nameA := Hash(name, HashNameA)
	nameB := Hash(name, HashNameB)
	for i := index; i < uint32(len(r.hashTable)); i++ {
		he := &r.hashTable[i]
		if he.blockIndex == 0xffffffff {
			break
		}
		if he.pathHashA == nameA && he.pathHashB == nameB {
			return he
		}
	}
	return nil
}

// ReadFile reads a file from within the MPQ file.
func (r *Reader) OpenFile(name string) io.Reader {
	he := r.findFile(name)
	if he == nil {
		return nil
	}
	be := r.blockTable[he.blockIndex]
	// log.Printf("blockEntry %#v", be)

	_, err := r.Seek(int64(r.userData.headerOfs+be.offset), 0)
	check(err)

	if be.flags&BlockFlagCompressed != 0 {
		comp := read8(r)
		if comp == 0x10 {
			return bzip2.NewReader(&io.LimitedReader{R: r, N: int64(be.size) - 1})
		}
		panic(fmt.Errorf("compressed file with unknown compression %x", comp))
	}
	panic("uncompressed not implemented")
}

// GetFileList returns a list of the files contained in the MPQ
// according to its "(listfile") metafile, if present.
func (r *Reader) GetFileList() []string {
	fr := r.OpenFile("(listfile)")
	if fr == nil {
		return nil
	}

	files := []string{}
	s := bufio.NewScanner(fr)
	for s.Scan() {
		files = append(files, s.Text())
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
	return files
}

func (r *Reader) readSection() {
	var buf [4]byte
	_, err := r.Read(buf[:])
	check(err)
	switch string(buf[:]) {
	case "MPQ\x1b": // user data
		r.readUserData()
	case "MPQ\x1a": // file header
		r.readHeader()
		r.readHashTable()
		r.readBlockTable()
		//r.readHET()
	default:
		panic(fmt.Errorf("section id: %q", buf))
	}
}

func (r *Reader) readHeaders() {
	r.readSection()
	if r.userData.headerOfs != 0 {
		_, err := r.Seek(int64(r.userData.headerOfs), 0)
		check(err)
		r.readSection()
	}
}

// NewReader reads the header of a file, returning an opened Reader.
func NewReader(f *os.File) *Reader {
	r := &Reader{File: f}
	r.readHeaders()
	return r
}
