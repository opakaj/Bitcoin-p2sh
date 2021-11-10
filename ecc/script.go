package ecc

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func p2pkhScript(h160 interface{}) *Script {
	//"Takes a hash160 and returns the p2pkh ScriptPubKey"
	return NewScript([]interface{}{118, 169, h160, 136, 172})
}

type Script struct {
	cmds []interface{}
}

func NewScript(cmds []interface{}) (s *Script) {
	cmds = nil
	s = new(Script)
	if &cmds == nil {
		s.cmds = []interface{}{}
	} else {
		s.cmds = cmds
	}
	return
}

func (s *Script) Repr() {
	var result []string
	var name string
	for _, cmd := range s.cmds {
		if reflect.TypeOf(cmd).Kind() == reflect.Int {
			if OPCODENAMES[cmd.(int)] != "" {
				name = OPCODENAMES[cmd.(int)] //didnt use get function
			} else {
				name = fmt.Sprintf("OP_[%d]", cmd)
			}
			result = append(result, name)
		} else {
			result = append(result, strconv.FormatInt(cmd.(int64), 16))
		}
	}
	strings.Join(result, " ")
}

func (S *Script) parse(s []byte) *Script {
	length := readVarint(s)
	var cmds []interface{}
	count := 0
	for int64(count) < length {
		//get the current byte
		var byt bytes.Buffer
		byt.Write(s)
		current, _ := byt.ReadByte()
		//increment the bytes we've read
		count += 1
		//convert current byte to integer
		currentByte := int(current)
		//if the current byte is between 1 and 75 inclusive
		if currentByte >= 1 && currentByte <= 75 {
			n := currentByte
			//add the next n bytes as an cmd
			x, _ := byt.ReadBytes(byte(n))
			cmds = append(cmds, x)
			//increase the count by n
			count += n
		} else if currentByte == 76 {
			//op_pushdata1
			x, _ := byt.ReadBytes(1)
			dataLength := littleEndianToInt(x)
			y, _ := byt.ReadBytes(byte(dataLength))
			cmds = append(cmds, y)
			count += int(dataLength) + 1
		} else if currentByte == 77 {
			//op_pushdata2
			x, _ := byt.ReadBytes(2)
			dataLength := littleEndianToInt(x)
			y, _ := byt.ReadBytes(byte(dataLength))
			cmds = append(cmds, y)
			count += int(dataLength) + 2
		} else {
			op_code := currentByte
			cmds = append(cmds, op_code)
		}
	}
	if count != int(length) {
		panic(fmt.Errorf("SyntaxError: %v", "parsing script failed"))
	}
	return NewScript(cmds)
}

func (s *Script) rawSerialize() []byte {
	// initialize what we'll send back
	result := []byte("")
	//go through each cmd
	for _, cmd := range s.cmds {
		//if the cmd is an integer, it's an opcode
		if reflect.TypeOf(cmd).Kind() == reflect.Int {
			//turn the cmd into a single byte integer using int_to_little_endian
			result = append(result, intToLittleEndian(cmd.(int), 1)...)
		} else {
			//otherwise, this is an element
			//get the length in bytes
			length := len(cmd)
			//for large lengths, we have to use a pushdata opcode
			if length < 75 {
				result = append(result, intToLittleEndian(length, 1)...)
			} else if length > 75 && length < 256 {
				//76 is pushdata1
				result = append(result, intToLittleEndian(76, 1)...)
				result = append(result, intToLittleEndian(length, 1)...)
			} else if length >= 256 && length <= 520 {
				//77 is pushdata2
				result = append(result, intToLittleEndian(77, 1)...)
				result = append(result, intToLittleEndian(length, 2)...)
			} else {
				panic(fmt.Errorf("ValueError: %v", "too long an cmd"))
			}
			result = append(result, cmd.(byte))
		}
	}
	return result
}

func (s *Script) serialize() []byte {
	result := s.rawSerialize()
	total := len(result)
	return append(encodeVarint(total), result...)
}

func (s *Script) evaluate(z interface{}) bool {
	cmds := s.cmds[:]
	var stack []byte
	var altstack []byte
	for len(cmds) > 0 {
		cmd := func(s *[]interface{}, i int) byte {
			popped := (*s)[i]
			*s = append((*s)[:i], (*s)[i+1:]...)
			return popped.(byte)
		}(&cmds, 0)
		if reflect.TypeOf(cmd).Kind() == reflect.Int {
			operation := OPCODEFUNCTIONS[int(cmd)]
			if func() int {
				for i, v := range [2]int{99, 100} {
					if byte(v) == cmd {
						return i
					}
				}
				return -1
			}() != -1 {
				if !operation(stack, cmds) {
					log.Print("bad op: {}", OPCODENAMES[int(cmd)])
					return false
				}
			} else if func() int {
				for i, v := range [2]int{107, 108} {
					if byte(v) == cmd {
						return i
					}
				}
				return -1
			}() != -1 {
				if !operation(stack, altstack) {
					log.Print("bad op: {}", OPCODENAMES[int(cmd)])
					return false
				}
			} else if func() int {
				for i, v := range [4]int{172, 173, 174, 175} {
					if byte(v) == cmd {
						return i
					}
				}
				return -1
			}() != -1 {
				if !operation(stack, z) {
					log.Print("bad op: {}", OPCODENAMES[int(cmd)])
					return false
				}
			} else if !operation(stack) {
				log.Print("bad op: {}", OPCODENAMES[int(cmd)])
				return false
			}
		} else {
			stack = append(stack, cmd)
		}
	}
	if len(stack) == 0 {
		return false
	}
	if reflect.DeepEqual(func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack), []byte("")) {
		return false
	}
	return true
}

func (s *Script) isP2pkhScriptPubkey() bool {
	//"Returns whether this follows the
	//OP_DUP OP_HASH160 <20 byte hash> OP_EQUALVERIFY OP_CHECKSIG pattern."
	return len(s.cmds) == 5 && s.cmds[0] == 118 &&
		(s.cmds[0] == 118 && s.cmds[1] == 169) &&
		(s.cmds[0] == 118 && s.cmds[1] == 169 && reflect.TypeOf(s.cmds[2]).Kind() == reflect.Slice) &&
		(s.cmds[0] == 118 && s.cmds[1] == 169 && reflect.TypeOf(s.cmds[2]).Kind() == reflect.Slice && len(s.cmds[2].([]byte)) == 20) &&
		(s.cmds[0] == 118 && s.cmds[1] == 169 && reflect.TypeOf(s.cmds[2]).Kind() == reflect.Slice && len(s.cmds[2].([]byte)) == 20 && s.cmds[3] == 136) &&
		(s.cmds[0] == 118 && s.cmds[1] == 169 && reflect.TypeOf(s.cmds[2]).Kind() == reflect.Slice && len(s.cmds[2].([]byte)) == 20 && s.cmds[3] == 136 && s.cmds[4] == 172)
}

func (s *Script) isP2shScriptPubkey() bool {
	//"Returns whether this follows the
	//OP_HASH160 <20 byte hash> OP_EQUAL pattern."
	return len(s.cmds) == 3 && s.cmds[0] == 169 &&
		(s.cmds[0] == 169 && reflect.TypeOf(s.cmds[1]).Kind() == reflect.Slice) &&
		(s.cmds[0] == 169 && reflect.TypeOf(s.cmds[1]).Kind() == reflect.Slice && len(s.cmds[1].([]byte)) == 20) &&
		(s.cmds[0] == 169 && reflect.TypeOf(s.cmds[1]).Kind() == reflect.Slice && len(s.cmds[1].([]byte)) == 20 && s.cmds[2] == 135)
}
