package ecc

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"github.com/btcsuite/btcutil"
)

var SIGHASHALL = 1
var SIGHASHNONE = 2
var SIGHASHSINGLE = 3
var BASE58ALPHABET = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func hash160(s string) string {
	return hex.EncodeToString(btcutil.Hash160([]byte(s)))
}

func hash256(s string) string {
	//two rounds of sha256
	hash := sha256.Sum256([]byte(s))
	hash = sha256.Sum256(hash[:])
	return hex.EncodeToString(hash[:]) //convet to [] by slicicng it
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func encodeBase58(s string) string {
	count := 0
	for _, c := range s {
		if c == 0 {
			count += 1
		} else {
			break
		}
	}
	num := binary.BigEndian.Uint32([]byte(s)) //bytes to int
	prefix := strings.Repeat("1", count)
	result := ""
	for num > 0 {
		_, mod := divmod(int64(num), 58)
		result = string(BASE58ALPHABET[mod]) + result
	}
	return prefix + result
}

func encodeBase58Checksum(b string) string {
	return encodeBase58(b + hash256(b)[:4])
}

func decodeBase58(s string) string {
	num := 0
	for _, c := range s {
		num *= 58
		num += func() int {
			for i, val := range BASE58ALPHABET {
				if val == c {
					return i
				}
			}
			panic(errors.New("ValueError: element not found"))
		}()
	}
	combined := make([]byte, 25)
	binary.BigEndian.PutUint64(combined, uint64(num)) //int to bytes
	checksum := combined[len(combined)-4:]
	if hash256(string(combined[:len(combined)-4]))[:4] != string(checksum) {
		panic(
			fmt.Errorf("bad address: %s %s", checksum, hash256(string(combined[:len(combined)-4]))[:4]),
		)
	}
	return string(combined[1 : len(combined)-4])
}

func littleEndianToInt(b []byte) int64 {
	//little_endian_to_int takes byte sequence as a little-endian number.
	//Returns an integer

	return int64(binary.BigEndian.Uint32(b)) //bytes to int
}

func intToLittleEndian(n, length int) []byte {
	//endian_to_little_endian takes an integer and returns the little-endian
	//byte sequence of length"

	x := make([]byte, length)
	binary.LittleEndian.PutUint64(x, uint64(n)) //int to bytes
	return x
}

func readVarint(s []byte) int64 {
	//read_varint reads a variable integer from a stream
	var byt bytes.Buffer
	byt.Write(s)

	i, _ := byt.ReadByte()

	if i == 253 {
		// 0xfd means the next two bytes are the number
		x, _ := byt.ReadBytes(2)
		return littleEndianToInt(x)
	} else if i == 254 {
		// 0xfe means the next two bytes are the number
		x, _ := byt.ReadBytes(4)
		return littleEndianToInt(x)
	} else if i == 255 {
		// 0xfe means the next two bytes are the number
		x, _ := byt.ReadBytes(8)
		return littleEndianToInt(x)
	} else {
		//anything else is just the integer
		return int64(i)
	}

}

func encodeVarint(i int) []byte {
	//encodes an integer as a varint
	if i < 253 {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(i)) //not so sure big or little
		return bs
	} else if i < 65536 {
		return append([]byte("\u00fd"), intToLittleEndian(i, 2)...)
	} else if i < 4294967296 {
		return append([]byte("\u00fe"), intToLittleEndian(i, 4)...)
	} else if i < 18446744073709551616 {
		return append([]byte("\u00ff"), intToLittleEndian(i, 8)...)
	} else {
		panic(fmt.Errorf("ValueError: %v", "integer too large: %d", i))
	}
}
