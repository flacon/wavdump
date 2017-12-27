package main

import (
	"fmt"
	"io"
	"os"
)

var data = make([]byte, 4)

func help() {
	fmt.Println("wavdump -Prints information from the header of the WAV file")
	fmt.Println("")
	fmt.Println("Usage: wavdump INPUTFILE")
}

func readTag(file *os.File) string {
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return string(data)
}

func readUInt16(file *os.File) uint16 {
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return uint16(data[1])<<8 +
		uint16(data[0])
}

func readUInt32(file *os.File) uint32 {
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return uint32(data[3])<<24 +
		uint32(data[2])<<16 +
		uint32(data[1])<<8 +
		uint32(data[0])
}

func printUInt32(value uint32, desc string) {
	hex := fmt.Sprintf("%0.8X", value)
	dec := fmt.Sprintf("%d", value)
	fmt.Printf("%14s %10s     %s\n", hex, dec, desc)
}

func printUInt16(value uint16, desc string) {
	hex := fmt.Sprintf("%0.4X", value)
	dec := fmt.Sprintf("%d", value)
	fmt.Printf("%14s %10s     %s\n", hex, dec, desc)
}

func read(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	tag := readTag(file)
	chunkSize := readUInt32(file)
	format := readTag(file)

	fmt.Printf("%s\n", tag)
	printUInt32(chunkSize, "ChunkSize")
	fmt.Printf("%14s  %10s    Format\n", format, "")

	for {
		subchunkID := readTag(file)
		subchunkSize := readUInt32(file)
		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("\n%s\n", subchunkID)
		printUInt32(subchunkSize, "SubChunkSize")

		if subchunkID == "fmt " {
			audioFormat := readUInt16(file)
			numChannels := readUInt16(file)
			sampleRate := readUInt32(file)
			byteRate := readUInt32(file)
			blockAlign := readUInt16(file)
			bitsPerSample := readUInt16(file)

			printUInt16(audioFormat, "AudioFormat")
			printUInt16(numChannels, "NumChannels")
			printUInt32(sampleRate, "SampleRate")
			printUInt32(byteRate, "ByteRate")
			printUInt16(blockAlign, "BlockAlign")
			printUInt16(bitsPerSample, "BitsPerSample")

		}

		if subchunkID == "data" {
			return nil
		}

		pos, err = file.Seek(pos+int64(subchunkSize), io.SeekStart)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}

}

func main() {
	if len(os.Args) != 2 {
		help()
		os.Exit(1)
	}

	err := read(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
