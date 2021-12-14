package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Line struct {
	address     int64
	data        []byte
	value       interface{}
	description string
}

type Parser struct {
	file  *os.File
	lines []Line
}

func NewParser(file *os.File) (Parser, error) {
	res := Parser{
		file:  file,
		lines: []Line{},
	}

	return res, nil
}

func (p *Parser) parse(fileName string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	p.file, err = os.Open(fileName)
	if err != nil {
		return err
	}
	defer p.file.Close()

	magic := readBytes(p.file, 4)
	_, err = p.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	switch string(magic) {

	case "RIFF":
		p.readWav()

	case "riff":
		p.readWave64()

	default:
		return fmt.Errorf("file \"%s\" is corrupted or has an unsupported format", fileName)
	}

	return nil
}

func filePos(file *os.File) int64 {
	res, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	return res
}

func readBytes(file *os.File, size uint) []byte {
	var res = make([]byte, size)
	_, err := file.Read(res)
	if err != nil {
		panic(err)
	}
	return res

}

func (p *Parser) addEmptyLine() {
	p.lines = append(p.lines, Line{})
}

func (p *Parser) readFourCC(description string) string {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, 4)
	res.value = string(res.data)
	p.lines = append(p.lines, res)
	return string(res.data)
}

func (p *Parser) readUInt16(description string) {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, 2)
	res.value = binary.LittleEndian.Uint16(res.data)
	p.lines = append(p.lines, res)
}

func (p *Parser) readUInt32(description string) uint32 {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, 4)
	res.value = binary.LittleEndian.Uint32(res.data)
	p.lines = append(p.lines, res)
	return binary.LittleEndian.Uint32(res.data)
}

func (p *Parser) readUInt64(description string) uint64 {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, 8)
	res.value = binary.LittleEndian.Uint64(res.data)
	p.lines = append(p.lines, res)
	return binary.LittleEndian.Uint64(res.data)
}

func (p *Parser) readBytes(size uint, description string) []byte {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, size)
	res.value = nil
	p.lines = append(p.lines, res)
	return res.data
}
