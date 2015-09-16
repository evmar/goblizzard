// Package blizzval decodes the Blizzard tagged value structure.
package blizzval

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
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

func read32(r io.Reader) (val uint32) {
	check(binary.Read(r, binary.LittleEndian, &val))
	return
}

type Value interface{}

// ReadVarInt reads the blizzval variable-length integer format.
func ReadVarInt(r io.Reader) int64 {
	var val int64
	for ofs := uint(0); true; ofs += 7 {
		b := read8(r)
		val |= (int64(b) & 0x7f) << ofs
		if b&0x80 == 0 {
			break
		}
	}
	if val&1 != 0 {
		return -(val >> 1)
	} else {
		return val >> 1
	}
}

// Read reads an encoded value from an io.Reader.
func Read(r io.Reader) Value {
	b := read8(r)
	switch b {
	case 0x0:
		size := ReadVarInt(r)
		array := make([]Value, size)
		for i := 0; i < int(size); i++ {
			array[i] = Read(r)
		}
		return array
	case 0x2:
		len := ReadVarInt(r)
		buf := make([]byte, len)
		_, err := r.Read(buf)
		if err != nil {
			panic(err)
		}
		return string(buf)
	case 0x4:
		if read8(r) == 0 {
			return nil
		}
		return Read(r)
	case 0x5:
		size := ReadVarInt(r)
		dict := map[int]Value{}
		for i := 0; i < int(size); i++ {
			key := ReadVarInt(r)
			val := Read(r)
			dict[int(key)] = val
		}
		return dict
	case 0x6:
		return read8(r)
	case 0x7:
		return read32(r)
	case 0x9:
		return ReadVarInt(r)
	default:
		panic(fmt.Errorf("tracker event %x", b))
	}

	return nil
}

func printIndent(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Printf("  ")
	}
}

func printVal(e Value, indent int) {
	switch e := e.(type) {
	case map[int]Value:
		keys := []int{}
		for k, _ := range e {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		fmt.Printf("{\n")
		for _, k := range keys {
			printIndent(indent + 1)
			fmt.Printf("%d: ", k)
			printVal(e[k], indent+1)
			fmt.Printf("\n")
		}
		printIndent(indent)
		fmt.Printf("}")
	case []Value:
		fmt.Printf("[\n")
		for _, val := range e {
			printIndent(indent + 1)
			printVal(val, indent+1)
			fmt.Printf("\n")
		}
	case string:
		fmt.Printf("%q", e)
	case uint8:
		fmt.Printf("%#x", e)
	case uint32:
		fmt.Printf("%#x", e)
	case int64:
		fmt.Printf("%#x", e)
	case nil:
		fmt.Printf("nil")
	default:
		panic(fmt.Errorf("unhandled %T", e))
	}
}

func Dump(v Value) {
	printVal(v, 0)
}
