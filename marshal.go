package hyperminhash

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const MarshalVersion = 1

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (sk *Sketch) MarshalBinary() ([]byte, error) {
	b := &bytes.Buffer{}
	b.Grow(8 + len(sk.reg)*2)
	b.WriteByte(MarshalVersion)
	b.WriteByte(p)
	b.WriteByte(q)
	b.WriteByte(r)
	binary.Write(b, binary.BigEndian, int32(m))
	for _, d := range sk.reg {
		binary.Write(b, binary.BigEndian, uint16(d))
	}
	return b.Bytes(), nil
}

// ErrorTooShort is an error that UnmarshalBinary try to parse too short
// binary.
var ErrorTooShort = errors.New("too short binary")

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (sk *Sketch) UnmarshalBinary(b []byte) error {
	if len(b) < 8 {
		return ErrorTooShort
	}

	if b[0] != MarshalVersion {
		return fmt.Errorf("unsupported version: %d", b[0])
	}
	if b[1] != p || b[2] != q || b[3] != r {
		return fmt.Errorf("unmatch parameters: want=%d,%d,%d got=%d,%d,%d", p, q, r, b[1], b[2], b[3])
	}

	r := bytes.NewReader(b[4:])
	var dlen int32
	err := binary.Read(r, binary.BigEndian, &dlen)
	if err != nil {
		return err
	}
	if dlen != int32(m) {
		return fmt.Errorf("unepected data lenghth: want=%d got=%d", m, dlen)
	}

	for i := 0; i < int(dlen); i++ {
		var d uint16
		err := binary.Read(r, binary.BigEndian, &d)
		if err != nil {
			return err
		}
		sk.reg[i] = register(d)
	}
	return nil
}
