package justdb

type KeyDirEntry struct {
	FileId        uint32
	ValueSize     uint32
	ValuePosition uint32
	Timestamp     uint32
}
