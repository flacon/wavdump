package main

import (
	"fmt"
)

func (p *Parser) readWave64() {
	p.readGUID("riff GUID")
	p.readUInt64("File size")
	p.readGUID("Format GUID")

	for i := 0; i < 100; i++ {
		p.addEmptyLine()
		chunkId := p.readGUID("Chunk GUID")

		switch string(chunkId[0:4]) {
		case "fmt ":
			p.readWave64FmtChunk()

		case "data":
			p.readUInt64("Data chunk size")
			return

		default:
			chunkSize := int(p.readUInt64("Chunk size"))
			p.readCustomChunk(chunkSize)
		}
	}

	panic(fmt.Errorf("data chunk not found"))
}

func (p *Parser) readGUID(description string) string {
	res := Line{}
	res.address = filePos(p.file)
	res.description = description
	res.data = readBytes(p.file, 16)
	v := fmt.Sprintf("%s...", res.data[0:4])
	res.value = v
	p.lines = append(p.lines, res)
	return v
}

func (p *Parser) readWave64FmtChunk() {
	chunkSize := p.readUInt64("Chunk size")
	p.readUInt16("Format code")
	p.readUInt16("Number of channels")
	p.readUInt32("Sample Rate")
	p.readUInt32("Byte Rate")
	p.readUInt16("Block Align")
	p.readUInt16("Bits Per Sample")

	if chunkSize -16 - 8 > 16 { // - sizeof(GUID) - sizeof(chunkSize)
		p.readUInt16("Size of the extension")
		p.readUInt16("Number of valid bits")
		p.readUInt32("Speaker position mask")

		p.readBytes(4, "⎫ GUID, including the")
		p.readBytes(4, "⎬ data format code")
		p.readBytes(4, "⎪ 16 bytes")
		p.readBytes(4, "⎭")
	}
}
