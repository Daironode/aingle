
package utils

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"

	" github.com/Daironode/aingle/common/serialization"
)

const (
	COMPRESS_TYPE_ZLIB = iota
)

const (
	DEFAULT_COMPRESS_TYPE         = COMPRESS_TYPE_ZLIB
	EXPORT_BLOCK_METADATA_LEN     = 256
	EXPORT_BLOCK_METADATA_VERSION = 1
)

type ExportBlockMetadata struct {
	Version          byte
	CompressType     byte
	StartBlockHeight uint32
	EndBlockHeight   uint32
}

func NewExportBlockMetadata() *ExportBlockMetadata {
	return &ExportBlockMetadata{
		Version:      EXPORT_BLOCK_METADATA_VERSION,
		CompressType: DEFAULT_COMPRESS_TYPE,
	}
}

func (this *ExportBlockMetadata) Serialize(w io.Writer) error {
	metadata := make([]byte, EXPORT_BLOCK_METADATA_LEN, EXPORT_BLOCK_METADATA_LEN)
	buf := bytes.NewBuffer(nil)
	err := serialization.WriteByte(buf, this.Version)
	if err != nil {
		return err
	}
	err = serialization.WriteByte(buf, this.CompressType)
	if err != nil {
		return err
	}
	err = serialization.WriteUint32(buf, this.StartBlockHeight)
	if err != nil {
		return err
	}
	err = serialization.WriteUint32(buf, this.EndBlockHeight)
	if err != nil {
		return err
	}
	data := buf.Bytes()
	if len(data) > EXPORT_BLOCK_METADATA_LEN {
		return fmt.Errorf("metata len size larger than %d", EXPORT_BLOCK_METADATA_LEN)
	}
	copy(metadata, data)
	_, err = w.Write(metadata)
	return err
}

func (this *ExportBlockMetadata) Deserialize(r io.Reader) error {
	metadata := make([]byte, EXPORT_BLOCK_METADATA_LEN, EXPORT_BLOCK_METADATA_LEN)
	_, err := io.ReadFull(r, metadata)
	if err != nil {
		return err
	}
	if metadata[0] != EXPORT_BLOCK_METADATA_VERSION {
		return fmt.Errorf("version unmatch")
	}
	reader := bytes.NewBuffer(metadata)
	ver, err := serialization.ReadByte(reader)
	if err != nil {
		return err
	}
	this.Version = ver
	compressType, err := serialization.ReadByte(reader)
	if err != nil {
		return err
	}
	this.CompressType = compressType
	height, err := serialization.ReadUint32(reader)
	if err != nil {
		return err
	}
	this.StartBlockHeight = height
	height, err = serialization.ReadUint32(reader)
	if err != nil {
		return err
	}
	this.EndBlockHeight = height
	return nil
}

func CompressBlockData(data []byte, compressType byte) ([]byte, error) {
	switch compressType {
	case COMPRESS_TYPE_ZLIB:
		return ZLibCompress(data)
	default:
		return nil, fmt.Errorf("unknown compress type")
	}
}

func DecompressBlockData(data []byte, compressType byte) ([]byte, error) {
	switch compressType {
	case COMPRESS_TYPE_ZLIB:
		return ZLibDecompress(data)
	default:
		return nil, fmt.Errorf("unknown compress type")
	}
}

func ZLibCompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	zlibWriter := zlib.NewWriter(buf)
	_, err := zlibWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("zlibWriter.Write error %s", err)
	}
	zlibWriter.Close()
	return buf.Bytes(), nil
}

func ZLibDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	zlibReader, err := zlib.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("zlib.NewReader error %s", err)
	}
	defer zlibReader.Close()

	return ioutil.ReadAll(zlibReader)
}
