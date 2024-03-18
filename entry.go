package justdb

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

type Entry struct {
	Checksum  uint32
	Timestamp uint32
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

func (e *Entry) SetChecksum() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, e.Timestamp)
	binary.Write(buf, binary.LittleEndian, e.KeySize)
	binary.Write(buf, binary.LittleEndian, e.ValueSize)
	buf.Write(e.Key)
	buf.Write(e.Value)

	e.Checksum = crc32.ChecksumIEEE(buf.Bytes())
}

func (e *Entry) VerifyChecksum() bool {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, e.Timestamp)
	binary.Write(buf, binary.LittleEndian, e.KeySize)
	binary.Write(buf, binary.LittleEndian, e.ValueSize)
	buf.Write(e.Key)
	buf.Write(e.Value)

	checksum := crc32.ChecksumIEEE(buf.Bytes())

	return checksum == e.Checksum
}
