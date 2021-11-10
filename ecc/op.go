package ecc

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"math"
	"reflect"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

func encodeNum(num int) []byte {
	if reflect.DeepEqual(num, 0) {
		return []byte("")
	}
	absNum := math.Abs(float64(num))
	negative := int(num) < 0
	var result []byte
	for absNum != 0 {
		result = append(result, byte(int(absNum)&255))
		int(absNum) >>= 8
	}
	//if the top bit is set,
	//for negative numbers we ensure that the top bit is set
	//for positive numbers we ensure that the top bit is not set
	if result[len(result)-1]&128 != 0 {
		if negative {
			result = append(result, 128)
		} else {
			result = append(result, 0)
		}
	} else if negative {
		result[len(result)-1] |= 128
	}
	return []byte(result)
}

func reverse(str string) string {
	output := ""
	for _, char := range str {
		output = string(char) + output
	}
	return output
}

func decodeNum(element string) int {
	var negative bool
	var result int
	if reflect.DeepEqual(element, []byte("")) {
		return 0
	}
	bigEndian := reverse(element)
	if bigEndian[0]&128 != 0 {
		negative = true
		result = int(bigEndian[0]) & 127
	} else {
		negative = false
		result = int(bigEndian[0])
	}
	for _, c := range bigEndian[1:] {
		result <<= 8
		result += int(c)
	}
	if negative {
		return -result
	} else {
		return result
	}
}

func op0(stack []byte) bool {
	stack = append(stack, encodeNum(0)...)
	return true
}

func op1negate(stack []byte) bool {
	stack = append(stack, encodeNum(-1)...)
	return true
}

func op1(stack []byte) bool {
	stack = append(stack, encodeNum(1)...)
	return true
}

func op2(stack []byte) bool {
	stack = append(stack, encodeNum(2)...)
	return true
}

func op3(stack []byte) bool {
	stack = append(stack, encodeNum(3)...)
	return true
}

func op4(stack []byte) bool {
	stack = append(stack, encodeNum(4)...)
	return true
}

func op5(stack []byte) bool {
	stack = append(stack, encodeNum(5)...)
	return true
}

func op6(stack []byte) bool {
	stack = append(stack, encodeNum(6)...)
	return true
}

func op7(stack []byte) bool {
	stack = append(stack, encodeNum(7)...)
	return true
}

func op8(stack []byte) bool {
	stack = append(stack, encodeNum(8)...)
	return true
}

func op9(stack []byte) bool {
	stack = append(stack, encodeNum(9)...)
	return true
}

func op10(stack []byte) bool {
	stack = append(stack, encodeNum(10)...)
	return true
}

func op11(stack []byte) bool {
	stack = append(stack, encodeNum(11)...)
	return true
}

func op12(stack []byte) bool {
	stack = append(stack, encodeNum(12)...)
	return true
}

func op13(stack []byte) bool {
	stack = append(stack, encodeNum(13)...)
	return true
}

func op14(stack []byte) bool {
	stack = append(stack, encodeNum(14)...)
	return true
}

func op15(stack []byte) bool {
	stack = append(stack, encodeNum(15)...)
	return true
}

func op16(stack []byte) bool {
	stack = append(stack, encodeNum(16)...)
	return true
}

func opNop(stack []byte) bool {
	return true
}

func opIf(stack []byte, items []byte) bool {
	if len(stack) < 1 {
		return false
	}
	var trueItems []byte
	var falseItems []byte
	currentArray := trueItems
	found := false
	numEndifsNeeded := 1
	for len(items) > 0 {
		item := func(s *[]byte, i int) byte {
			popped := (*s)[i]
			*s = append((*s)[:i], (*s)[i+1:]...)
			return popped
		}(&items, 0)
		if func() int {
			for i, v := range [2]int{99, 100} {
				if byte(v) == item {
					return i
				}
			}
			return -1
		}() != -1 {
			numEndifsNeeded += 1
			currentArray = append(currentArray, item)
		} else if numEndifsNeeded == 1 && reflect.DeepEqual(item, 103) {
			currentArray = falseItems
		} else if reflect.DeepEqual(item, 104) {
			if numEndifsNeeded == 1 {
				found = true
				break
			} else {
				numEndifsNeeded -= 1
				currentArray = append(currentArray, item)
			}
		} else {
			currentArray = append(currentArray, item)
		}
	}
	if !found {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	if decodeNum(element) == 0 {
		items[:0] = falseItems
	} else {
		items[:0] = trueItems
	}
	return true
}

func opNotIf(stack, items []byte) bool {
	if len(stack) < 1 {
		return false
	}
	trueItems := []interface{}{}
	falseItems := []interface{}{}
	currentArray := trueItems
	found := false
	numEndifsNeeded := 1
	for len(items) > 0 {
		item := func(s *[]byte, i int) byte {
			popped := (*s)[i]
			*s = append((*s)[:i], (*s)[i+1:]...)
			return popped
		}(&items, 0)
		if func() int {
			for i, v := range [2]int{99, 100} {
				if byte(v) == item {
					return i
				}
			}
			return -1
		}() != -1 {
			numEndifsNeeded += 1
			currentArray = append(currentArray, item)
		} else if numEndifsNeeded == 1 && reflect.DeepEqual(item, 103) {
			currentArray = falseItems
		} else if reflect.DeepEqual(item, 104) {
			if numEndifsNeeded == 1 {
				found = true
				break
			} else {
				numEndifsNeeded -= 1
				currentArray = append(currentArray, item)
			}
		} else {
			currentArray = append(currentArray, item)
		}
	}
	if !found {
		return false
	}
	element := func(s *[]byte) byte {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	if decodeNum(string(element)) == 0 {
		items[:0] = trueItems
	} else {
		items[:0] = falseItems
	}
	return true
}

func opVerify(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string { //pass by pointer so that it actually changes
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	if decodeNum(element) == 0 {
		return false
	}
	return true
}

func opReturn(stack []byte) bool {
	return false
}

func opToAltStack(stack []byte, altstack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	altstack = append(altstack, func(s *[]byte) byte {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack))
	return true
}

func opFromAltStack(stack []byte, altstack []byte) bool {
	if len(altstack) < 1 {
		return false
	}
	stack = append(stack, func(s *[]byte) byte {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&altstack))
	return true
}

func op2Drop(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	return true
}

func op2Dup(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-2:]...)
	return true
}

func op3Dup(stack []byte) bool {
	if len(stack) < 3 {
		return false
	}
	stack = append(stack, stack[len(stack)-3:]...)
	return true
}

func op2Over(stack []byte) bool {
	if len(stack) < 4 {
		return false
	}
	stack = append(stack, stack[len(stack)-4:len(stack)-2]...)
	return true
}

func op2Rot(stack []byte) bool {
	if len(stack) < 6 {
		return false
	}
	stack = append(stack, stack[len(stack)-6:len(stack)-4]...)
	return true
}

func op2Swap(stack []byte) bool {
	if len(stack) < 4 {
		return false
	}
	stack[len(stack)-4:] = append(stack[len(stack)-2:], stack[len(stack)-4:len(stack)-2]...)
	return true
}

func opIfDup(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	if decodeNum(string(stack[len(stack)-1])) != 0 {
		stack = append(stack, stack[len(stack)-1])
	}
	return true
}

func opDepth(stack []byte) bool {
	stack = append(stack, encodeNum(len(stack))...)
	return true
}

func opDrop(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	return true
}

func opDup(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = append(stack, stack[len(stack)-1])
	return true
}

func opNip(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	//copy(stack[len(stack)-2:], stack[len(stack)-1:])
	stack[len(stack)-2:] = stack[len(stack)-1:]

	return true
}

func opOver(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-2])
	return true
}

func opPick(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	n := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if len(stack) < n+1 {
		return false
	}
	stack = append(stack, stack[-n-1])
	return true
}

func opRoll(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	n := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if len(stack) < n+1 {
		return false
	}
	if n == 0 {
		return true
	}
	stack = append(stack, func(s *[]byte, i int) byte {
		popped := (*s)[i]
		*s = append((*s)[:i], (*s)[i+1:]...)
		return popped
	}(&stack, -n-1))
	return true
}

func opRot(stack []byte) bool {
	if len(stack) < 3 {
		return false
	}
	stack = append(stack, func(s *[]byte, i int) byte {
		popped := (*s)[i]
		*s = append((*s)[:i], (*s)[i+1:]...)
		return popped
	}(&stack, -3))
	return true
}

func opSwap(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, func(s *[]byte, i int) byte {
		popped := (*s)[i]
		*s = append((*s)[:i], (*s)[i+1:]...)
		return popped
	}(&stack, -2))
	return true
}

func opTuck(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-1])
	copy(stack[len(stack)-2+1:], stack[len(stack)-2:])
	stack[len(stack)-2] = stack[len(stack)-1]
	return true
}

func opSize(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = append(stack, encodeNum(len(string(stack[len(stack)-1])))...)
	return true
}

func opEqual(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	element2 := func(s *[]byte) interface{} {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	if reflect.DeepEqual(element1, element2) {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opEqualVerify(stack []byte) bool {
	return opEqual(stack) && opVerify(stack)
}

func op1Add(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	stack = append(stack, encodeNum(element+1)...)
	return true
}

func op1Sub(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	stack = append(stack, encodeNum(element-1)...)
	return true
}

func opNegate(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	stack = append(stack, encodeNum(-element)...)
	return true
}

func opAbs(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element < 0 {
		stack = append(stack, encodeNum(-element)...)
	} else {
		stack = append(stack, encodeNum(element)...)
	}
	return true
}

func opNot(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	if decodeNum(element) == 0 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func op0NotEqual(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	if decodeNum(element) == 0 {
		stack = append(stack, encodeNum(0)...)
	} else {
		stack = append(stack, encodeNum(1)...)
	}
	return true
}

func opAdd(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	stack = append(stack, encodeNum(element1+element2)...)
	return true
}

func opSub(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	stack = append(stack, encodeNum(element2-element1)...)
	return true
}

func opBoolAnd(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 && element2 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opBoolOr(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 || element2 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opNumEqual(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 == element2 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opNumEqualVerify(stack []byte) bool {
	return opNumEqual(stack) && opVerify(stack)
}

func opNumNotEqual(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 == element2 {
		stack = append(stack, encodeNum(0)...)
	} else {
		stack = append(stack, encodeNum(1)...)
	}
	return true
}

func opLessThan(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element2 < element1 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opGreaterThan(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element2 > element1 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opLessThanOrEqual(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element2 <= element1 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opGreaterThanOrEqual(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element2 >= element1 {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opMin(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 < element2 {
		stack = append(stack, encodeNum(element1)...)
	} else {
		stack = append(stack, encodeNum(element2)...)
	}
	return true
}

func opMax(stack []byte) bool {
	if len(stack) < 2 {
		return false
	}
	element1 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element2 := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element1 > element2 {
		stack = append(stack, encodeNum(element1)...)
	} else {
		stack = append(stack, encodeNum(element2)...)
	}
	return true
}

func opWithin(stack []byte) bool {
	if len(stack) < 3 {
		return false
	}
	maximum := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	minimum := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	element := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if element >= minimum && element < maximum {
		stack = append(stack, encodeNum(1)...)
	} else {
		stack = append(stack, encodeNum(0)...)
	}
	return true
}

func opRipemd160(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	stack = append(stack, func([]byte) []byte {
		first := sha256.Sum256([]byte(element))
		hasher := ripemd160.New()
		hasher.Write(first[:])
		hash := hasher.Sum(nil)
		return hash[:]
	}([]byte(element))...) //hashlib.new("ripemd160", element).digest())
	return true
}

func opSha1(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	stack = append(stack, func() []byte {
		sha1 := sha1.Sum([]byte(element))
		return sha1[:]
	}()...) //hashlib.sha1(element).digest())
	return true
}

func opSha256(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	stack = append(stack, func() []byte {
		sha256 := sha256.Sum256([]byte(element))
		return sha256[:]
	}()...)
	return true
}

func opHash160(stack string) {
	panic(errors.New("Exception"))
}

func opHash256(stack []byte) bool {
	if len(stack) < 1 {
		return false
	}
	element := func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack)
	stack = append(stack, hash256(element)...)
	return true
}

func opCheckSig(stack []byte, z interface{}) bool {
	panic(errors.New("Exception"))
}

func opCheckSigVerify(stack []byte, z interface{}) bool {
	return opCheckSig(stack, z) && opVerify(stack)
}

func opCheckMultiSig(stack []byte, z interface{}) bool {
	if len(stack) < 1 {
		return false
	}
	n := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if len(stack) < n+1 {
		return false
	}
	var secPubkeys []byte
	for i := 0; i < n; i++ {
		secPubkeys = append(secPubkeys, func(s *[]byte) byte {
			i := len(*s) - 1
			popped := (*s)[i]
			*s = (*s)[:i]
			return popped
		}(&stack))
	}
	m := decodeNum(func(s *[]byte) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return string(popped)
	}(&stack))
	if len(stack) < m+1 {
		return false
	}
	var derSignatures []byte
	for x := 0; x < m; x++ {
		derSignatures = append(derSignatures, func(s *[]byte) []byte {
			i := len(*s) - 1
			popped := (*s)[i]
			*s = (*s)[:i]
			return []byte(popped)
		}(&stack)[:len(stack)-1]...)
	}
	func(s *[]byte) byte {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = (*s)[:i]
		return popped
	}(&stack)
	func() bool {
		defer func() bool {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					if strings.HasPrefix(err.Error(), "ValueError") ||
						strings.HasPrefix(err.Error(), "SyntaxError") {
						return false
					}
				}
				panic(r)
			}
		}()
		points := func() (elts []interface{}) {
			for _, sec := range secPubkeys {
				s256 := new(S256Point)
				elts = append(elts, s256.parse(sec))
			}
			return
		}()
		sigs := func() (elts []interface{}) {
			for _, der := range derSignatures {
				s := new(Signature)
				elts = append(elts, s.parse(der))
			}
			return
		}()
		for _, sig := range sigs {
			if len(points) == 0 {
				return false
			}
			for len(points) != 0 {
				point := func(s *[]interface{}, i int) interface{} {
					popped := (*s)[i]
					*s = append((*s)[:i], (*s)[i+1:]...)
					return popped
				}(&points, 0)
				if point.verify(z, sig) {
					break
				}
			}
		}
		stack = append(stack, encodeNum(1)...)
	}()
	return true
}

func opCheckMultiSigVerify(stack []byte, z interface{}) bool {
	return opCheckMultiSig(stack, z) && opVerify(stack)
}

func opCheckLocktimeVerify(stack []byte, locktime int, sequence interface{}) bool {
	if reflect.DeepEqual(sequence, 4294967295) {
		return false
	}
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(string(stack[len(stack)-1]))
	if element < 0 {
		return false
	}
	if element < 500000000 && int(locktime) > 500000000 {
		return false
	}
	if locktime < element {
		return false
	}
	return true
}

func opCheckSequenceVerify(stack []byte, version int, sequence int) bool {
	if reflect.DeepEqual(int(sequence)&(1<<31), 1<<31) {
		return false
	}
	if len(stack) < 1 {
		return false
	}
	element := decodeNum(string(stack[len(stack)-1]))
	if element < 0 {
		return false
	}
	if element&(1<<31) == 1<<31 {
		if int(version) < 2 {
			return false
		} else if reflect.DeepEqual(int(sequence)&(1<<31), 1<<31) {
			return false
		} else if !reflect.DeepEqual(element&(1<<22), int(sequence)&(1<<22)) {
			return false
		} else if element&65535 > int(sequence)&65535 {
			return false
		}
	}
	return true
}

var OPCODEFUNCTIONS = map[int]interface{}{
	0:   op0,
	79:  op1negate,
	81:  op1,
	82:  op2,
	83:  op3,
	84:  op4,
	85:  op5,
	86:  op6,
	87:  op7,
	88:  op8,
	89:  op9,
	90:  op10,
	91:  op11,
	92:  op12,
	93:  op13,
	94:  op14,
	95:  op15,
	96:  op16,
	97:  opNop,
	99:  opIf,
	100: opNotIf,
	105: opVerify,
	106: opReturn,
	107: opToAltStack,
	108: opFromAltStack,
	109: op2Drop,
	110: op2Dup,
	111: op3Dup,
	112: op2Over,
	113: op2Rot,
	114: op2Swap,
	115: opIfDup,
	116: opDepth,
	117: opDrop,
	118: opDup,
	119: opNip,
	120: opOver,
	121: opPick,
	122: opRoll,
	123: opRot,
	124: opSwap,
	125: opTuck,
	130: opSize,
	135: opEqual,
	136: opEqualVerify,
	139: op1Add,
	140: op1Sub,
	143: opNegate,
	144: opAbs,
	145: opNot,
	146: op0NotEqual,
	147: opAdd,
	148: opSub,
	154: opBoolAnd,
	155: opBoolOr,
	156: opNumEqual,
	157: opNumEqualVerify,
	158: opNumNotEqual,
	159: opLessThan,
	160: opGreaterThan,
	161: opLessThanOrEqual,
	162: opGreaterThanOrEqual,
	163: opMin,
	164: opMax,
	165: opWithin,
	166: opRipemd160,
	167: opSha1,
	168: opSha256,
	169: opHash160,
	170: opHash256,
	172: opCheckSig,
	173: opCheckSigVerify,
	174: opCheckMultiSig,
	175: opCheckMultiSigVerify,
	176: opNop,
	177: opCheckLocktimeVerify,
	178: opCheckSequenceVerify,
	179: opNop,
	180: opNop,
	181: opNop,
	182: opNop,
	183: opNop,
	184: opNop,
	185: opNop,
}

var OPCODENAMES = map[int]string{
	0:   "OP0",
	76:  "OPPUSHDATA1",
	77:  "OPPUSHDATA2",
	78:  "OPPUSHDATA4",
	79:  "OP1NEGATE",
	81:  "OP1",
	82:  "OP2",
	83:  "OP3",
	84:  "OP4",
	85:  "OP5",
	86:  "OP6",
	87:  "OP7",
	88:  "OP8",
	89:  "OP9",
	90:  "OP10",
	91:  "OP11",
	92:  "OP12",
	93:  "OP13",
	94:  "OP14",
	95:  "OP15",
	96:  "OP16",
	97:  "OPNOP",
	99:  "OPIF",
	100: "OPNOTIF",
	103: "OPELSE",
	104: "OPENDIF",
	105: "OPVERIFY",
	106: "OPRETURN",
	107: "OPTOALTSTACK",
	108: "OPFROMALTSTACK",
	109: "OP2DROP",
	110: "OP2DUP",
	111: "OP3DUP",
	112: "OP2OVER",
	113: "OP2ROT",
	114: "OP2SWAP",
	115: "OPIFDUP",
	116: "OPDEPTH",
	117: "OPDROP",
	118: "OPDUP",
	119: "OPNIP",
	120: "OPOVER",
	121: "OPPICK",
	122: "OPROLL",
	123: "OPROT",
	124: "OPSWAP",
	125: "OPTUCK",
	130: "OPSIZE",
	135: "OPEQUAL",
	136: "OPEQUALVERIFY",
	139: "OP1ADD",
	140: "OP1SUB",
	143: "OPNEGATE",
	144: "OPABS",
	145: "OPNOT",
	146: "OP0NOTEQUAL",
	147: "OPADD",
	148: "OPSUB",
	154: "OPBOOLAND",
	155: "OPBOOLOR",
	156: "OPNUMEQUAL",
	157: "OPNUMEQUALVERIFY",
	158: "OPNUMNOTEQUAL",
	159: "OPLESSTHAN",
	160: "OPGREATERTHAN",
	161: "OPLESSTHANOREQUAL",
	162: "OPGREATERTHANOREQUAL",
	163: "OPMIN",
	164: "OPMAX",
	165: "OPWITHIN",
	166: "OPRIPEMD160",
	167: "OPSHA1",
	168: "OPSHA256",
	169: "OPHASH160",
	170: "OPHASH256",
	171: "OPCODESEPARATOR",
	172: "OPCHECKSIG",
	173: "OPCHECKSIGVERIFY",
	174: "OPCHECKMULTISIG",
	175: "OPCHECKMULTISIGVERIFY",
	176: "OPNOP1",
	177: "OPCHECKLOCKTIMEVERIFY",
	178: "OPCHECKSEQUENCEVERIFY",
	179: "OPNOP4",
	180: "OPNOP5",
	181: "OPNOP6",
	182: "OPNOP7",
	183: "OPNOP8",
	184: "OPNOP9",
	185: "OPNOP10",
}
