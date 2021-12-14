package main

import (
	"fmt"
)

func (p *Parser) readWav() {
	p.readFourCC("RIFF")
	p.readUInt32("File size - 8")
	p.readFourCC("Format")

	for i := 0; i < 100; i++ {
		p.addEmptyLine()
		chunkId := p.readFourCC("Chunk ID")

		switch chunkId {
		case "fmt ":
			p.readWavFmtChunk()

		case "data":
			p.readUInt32("Data chunk size")
			return

		default:
			chunkSize := int(p.readUInt32("Chunk size"))
			p.readCustomChunk(chunkSize)
		}
	}

	panic(fmt.Errorf("data chunk not found"))
}

func (p *Parser) readWavFmtChunk() {
	chunkSize := p.readUInt32("Chunk size")
	p.readUInt16("Format code")
	p.readUInt16("Number of channels")
	p.readUInt32("Sample Rate")
	p.readUInt32("Byte Rate")
	p.readUInt16("Block Align")
	p.readUInt16("Bits Per Sample")

	if chunkSize > 16 {
		p.readUInt16("Size of the extension")
		p.readUInt16("Number of valid bits")
		p.readUInt32("Speaker position mask")

		p.readBytes(4, "⎫ GUID, including the")
		p.readBytes(4, "⎬ data format code")
		p.readBytes(4, "⎪ 16 bytes")
		p.readBytes(4, "⎭")
	}
}

func (p *Parser) readCustomChunk(chunkSize int) {

	cnt := 0
	for i := 0; i < chunkSize; i += 4 {
		sz := 4
		if i+sz >= chunkSize {
			sz = chunkSize - i
		}

		cnt++
		p.readBytes(uint(sz), "⎪")
	}

	switch cnt {
	case 0:
		return

	case 1:
		p.lines[len(p.lines)-1].description = "Data"

	default:
		p.lines[len(p.lines)-cnt].description = "⎫ Data"
		p.lines[len(p.lines)-1].description = "⎭"

	}
}
