package main

import (
	"fmt"
	"io"
	"os"
)

func help() {
	fmt.Println("wavdump -Prints information from the header of the WAV file")
	fmt.Println("")
	fmt.Println("Usage: wavdump INPUTFILE")
}

func readArray(file *os.File, size uint) []byte {
	var data = make([]byte, size)
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return data
}

func readTag(file *os.File) ([]byte, string) {
	var data = make([]byte, 4)
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return data, string(data)
}

func readUInt16(file *os.File) ([]byte, uint16) {
	var data = make([]byte, 2)
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return data,
		uint16(data[1])<<8 +
			uint16(data[0])
}

func readUInt32(file *os.File) ([]byte, uint32) {
	var data = make([]byte, 4)
	_, err := file.Read(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return data,
		uint32(data[3])<<24 +
			uint32(data[2])<<16 +
			uint32(data[1])<<8 +
			uint32(data[0])

}

func printBytes(data []byte) string {
	res := ""
	for _, b := range data {
		res += fmt.Sprintf("%0.2X ", b)
	}
	return res
}

func print(col1 string, col2 interface{}, bytes []byte, comment string) {
	s := fmt.Sprintf("%-30s %16v ", col1, col2)
	fmt.Printf("%s %16s    %s\n", s, printBytes(bytes), comment)
}

func read(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data, tag := readTag(file)
	print(tag, "", data, "")

	data, chunkSize := readUInt32(file)
	print("ChunkSize", chunkSize, data, "")

	data, format := readTag(file)
	print("Format", format, data, "")

	for {
		fmt.Println()

		data, subchunkID := readTag(file)
		print(subchunkID, "", data, "Chunk ID")

		data, subchunkSize := readUInt32(file)
		print("    cksize", subchunkSize, data, "Chunk size")

		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		if subchunkID == "fmt " {
			data, audioFormat := readUInt16(file)
			print("    wFormatTag", audioFormat, data, "Format code")

			data, numChannels := readUInt16(file)
			print("    NumChannels", numChannels, data, "")

			data, sampleRate := readUInt32(file)
			print("    SampleRate", sampleRate, data, "")

			data, byteRate := readUInt32(file)
			print("    ByteRate", byteRate, data, "")

			data, blockAlign := readUInt16(file)
			print("    BlockAlign", blockAlign, data, "")

			data, bitsPerSample := readUInt16(file)
			print("    BitsPerSample", bitsPerSample, data, "")

			if subchunkSize > 16 {
				data, cbSize := readUInt16(file)
				print("      sbSize", cbSize, data, "Size of the extension")

				data, wValidBitsPerSample := readUInt16(file)
				print("      wValidBitsPerSample", wValidBitsPerSample, data, "Number of valid bits")

				data, dwChannelMask := readUInt32(file)
				print("      dwChannelMask", dwChannelMask, data, "Speaker position mask")

				subFormat := readArray(file, 16)
				print("      SubFormat", "", subFormat[0:4], "⎫  GUID, including the")
				print("", "", subFormat[4:8], "⎬  data format code")
				print("", "", subFormat[8:12], "⎥  16 bytes")
				print("", "", subFormat[12:16], "⎭")
			}
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
