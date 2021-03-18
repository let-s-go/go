package jsoniter

import (
	"encoding/json"
	"io"
	"math"
	"strconv"
	"unsafe"
)

func UseNumberAsString() {
	RegisterTypeDecoderFunc("string", func(ptr unsafe.Pointer, iter *Iterator) {
		next := iter.WhatIsNext()
		switch next {
		case NumberValue:
			var number json.Number
			iter.ReadVal(&number)
			*((*string)(ptr)) = string(number)
		case StringValue:
			*((*string)(ptr)) = iter.ReadString()
		default:
			if iter.tryReadNaN() {
				*((*string)(ptr)) = "NaN"
			} else {
				iter.ReportError("StringDecoder", "not number or string")
			}
		}
	})
}

func UseStringAsNumber() {
	RegisterTypeDecoderFunc("bool", func(ptr unsafe.Pointer, iter *Iterator) {
		if b, ok := iter.readStringAsBool(); ok {
			*((*bool)(ptr)) = b
		} else {
			*((*bool)(ptr)) = iter.ReadBool()
		}
	})
	RegisterTypeDecoderFunc("byte", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(8); ok {
			*((*byte)(ptr)) = byte(n)
		} else {
			*((*byte)(ptr)) = iter.readByte()
		}
	})
	RegisterTypeDecoderFunc("int", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(32); ok {
			*((*int)(ptr)) = int(n)
		} else {
			*((*int)(ptr)) = iter.ReadInt()
		}
	})
	RegisterTypeDecoderFunc("int8", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(8); ok {
			*((*int8)(ptr)) = int8(n)
		} else {
			*((*int8)(ptr)) = iter.ReadInt8()
		}
	})
	RegisterTypeDecoderFunc("int16", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(16); ok {
			*((*int16)(ptr)) = int16(n)
		} else {
			*((*int16)(ptr)) = iter.ReadInt16()
		}
	})
	RegisterTypeDecoderFunc("int32", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(32); ok {
			*((*int32)(ptr)) = int32(n)
		} else {
			*((*int32)(ptr)) = iter.ReadInt32()
		}
	})
	RegisterTypeDecoderFunc("int64", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsSigned(64); ok {
			*((*int64)(ptr)) = int64(n)
		} else {
			*((*int64)(ptr)) = iter.ReadInt64()
		}
	})
	RegisterTypeDecoderFunc("uint", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsUnSigned(32); ok {
			*((*uint)(ptr)) = uint(n)
		} else {
			*((*uint)(ptr)) = iter.ReadUint()
		}
	})
	RegisterTypeDecoderFunc("uint8", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsUnSigned(8); ok {
			*((*uint8)(ptr)) = uint8(n)
		} else {
			*((*uint8)(ptr)) = iter.ReadUint8()
		}
	})
	RegisterTypeDecoderFunc("uint16", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsUnSigned(16); ok {
			*((*uint16)(ptr)) = uint16(n)
		} else {
			*((*uint16)(ptr)) = iter.ReadUint16()
		}
	})
	RegisterTypeDecoderFunc("uint32", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsUnSigned(32); ok {
			*((*uint32)(ptr)) = uint32(n)
		} else {
			*((*uint32)(ptr)) = iter.ReadUint32()
		}
	})
	RegisterTypeDecoderFunc("uint64", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsUnSigned(64); ok {
			*((*uint64)(ptr)) = uint64(n)
		} else {
			*((*uint64)(ptr)) = iter.ReadUint64()
		}
	})
	RegisterTypeDecoderFunc("float32", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsFloat(32); ok {
			*((*float32)(ptr)) = float32(n)
		} else {
			if iter.tryReadNaN() {
				*((*float32)(ptr)) = float32(math.NaN())
			} else {
				*((*float32)(ptr)) = iter.ReadFloat32()
			}
		}
	})
	RegisterTypeDecoderFunc("float64", func(ptr unsafe.Pointer, iter *Iterator) {
		if n, ok := iter.readStringAsFloat(64); ok {
			*((*float64)(ptr)) = float64(n)
		} else {
			if iter.tryReadNaN() {
				*((*float64)(ptr)) = math.NaN()
			} else {
				*((*float64)(ptr)) = iter.ReadFloat64()
			}
		}
	})
}

func (iter *Iterator) readStringAsBool() (bool, bool) {
	if iter.WhatIsNext() == StringValue {
		s := iter.ReadString()
		if s == "true" {
			return true, true
		} else if s == "false" {
			return false, true
		} else {
			iter.ReportError("ParseBool", "cant't parse string to bool")
			return false, true
		}
	}
	return false, false
}

func (iter *Iterator) readStringAsSigned(bitSize int) (int64, bool) {
	next := iter.WhatIsNext()
	if next == StringValue {
		s := iter.ReadString()
		if iter.Error != nil && iter.Error != io.EOF {
			return 0, true
		}
		n, err := strconv.ParseInt(s, 10, bitSize)
		if err != nil {
			iter.ReportError("ParseInt", err.Error())
			return 0, true
		}
		return n, true
	} else if next == NilValue {
		if !iter.ReadNil() {
			iter.ReportError("ParseInt", "start with n but not null")
			return 0, true
		}
		return 0, true
	}
	return 0, false
}

func (iter *Iterator) readStringAsUnSigned(bitSize int) (uint64, bool) {
	next := iter.WhatIsNext()
	if next == StringValue {
		s := iter.ReadString()
		if iter.Error != nil && iter.Error != io.EOF {
			return 0, true
		}
		n, err := strconv.ParseUint(s, 10, bitSize)
		if err != nil {
			iter.ReportError("ParseUint", err.Error())
			return 0, true
		}
		return n, true
	} else if next == NilValue {
		if !iter.ReadNil() {
			iter.ReportError("ParseInt", "start with n but not null")
			return 0, true
		}
		return 0, true
	}
	return 0, false
}

func (iter *Iterator) readStringAsFloat(bitSize int) (float64, bool) {
	next := iter.WhatIsNext()
	if next == StringValue {
		s := iter.ReadString()
		if iter.Error != nil && iter.Error != io.EOF {
			return 0, true
		}
		n, err := strconv.ParseFloat(s, bitSize)
		if err != nil {
			iter.ReportError("ParseFloat", err.Error())
			return 0, true
		}
		return n, true
	} else if next == NilValue {
		if !iter.ReadNil() {
			iter.ReportError("ParseFloat", "start with n but not null")
			return 0, true
		}
		return 0, true
	}
	return 0, false
}

// ReadNil reads a json object as nil and
// returns whether it's a nil or not
func (iter *Iterator) tryReadNaN() bool {
	if iter.nextToken() == 'N' {
		a := iter.readByte()
		b := iter.readByte()
		if a == 'a' && b == 'N' {
			return true
		}
		iter.unreadByte()
		iter.unreadByte()
	}
	iter.unreadByte()
	return false
}
