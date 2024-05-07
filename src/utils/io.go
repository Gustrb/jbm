package utils

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type BigEndianReader struct {
	reader *bytes.Reader
}

func NewBigEndianReaderFromReader(reader *bytes.Reader) *BigEndianReader {
	return &BigEndianReader{
		reader: reader,
	}
}

func (ber *BigEndianReader) ReadUint8() (uint8, error) {
	byte := make([]byte, 1)

	_, err := ber.reader.Read(byte)
	if err != nil {
		return 0, err
	}

	return byte[0], nil
}

func (ber *BigEndianReader) ReadUint16() (uint16, error) {
	bytes := make([]byte, 2)

	_, err := ber.reader.Read(bytes)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(bytes), nil
}

func (ber *BigEndianReader) ReadUint32() (uint32, error) {
	bytes := make([]byte, 4)

	_, err := ber.reader.Read(bytes)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(bytes), nil
}

func (ber *BigEndianReader) ReadUint64() (uint64, error) {
	bytes := make([]byte, 8)

	_, err := ber.reader.Read(bytes)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(bytes), nil
}

func (ber *BigEndianReader) ReadBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)

	_, err := ber.reader.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ReadFileContent(filepath string) ([]byte, error) {
	fptr, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer fptr.Close()

	fcontent, err := io.ReadAll(fptr)

	if err != nil {
		return nil, err
	}

	return fcontent, nil
}
