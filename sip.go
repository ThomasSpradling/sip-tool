package main

import (
	"encoding/binary"
	"errors"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Compresses file contents and places in `compressed`. Also creates
// a `decodeMap` which encodes the Huffman tree in a decodable format.
func compress(contents string) (decodeMap []byte, compressed []byte) {
	return nil, nil
}

// Creates .sipped file, and compresses directories into it.
func sip(directories []string, name string) error {
	if len(name) == 0 {
		name = "archive"
	}

	var stack []string

	for _, dir := range directories {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			return errors.New("Not all files or directories could be found.")
		}

		stack = append(stack, dir)
	}

	f, err := os.Create(name + ".sipped")
	check(err)
	defer f.Close()

	for len(stack) > 0 {
		// Grab current directory to consider.
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		fileInfo, err := os.Stat(curr)
		check(err)

		// If it is a directory, add subdirectories of folder to stack.
		if fileInfo.IsDir() {
			subdirectories, err := os.ReadDir(curr)
			check(err)

			for _, dir := range subdirectories {
				stack = append(stack, curr+"/"+dir.Name())
			}
		} else {
			// Otherwise, compress and add to .sipped.

			buf, err := os.ReadFile(curr)
			check(err)

			// If filename takes up more than 256 bytes,
			// we cannot store with current format.
			if len(curr) > 256 {
				panic("Filename too long!")
			}

			decodeMap, compressed := compress(string(buf))
			var filenameSize uint8 = uint8(len(curr))
			var compressionSize uint16 = uint16(len(compressed))
			var decodeMapSize uint8 = uint8(len(decodeMap))

			// Write a header with format:
			// filename_size (bytes) | filename |
			// compression_size (bytes) |
			// decode_map_size (bytes) | decode_map_bits

			binary.Write(f, binary.LittleEndian, filenameSize)
			f.WriteString(curr)

			binary.Write(f, binary.LittleEndian, compressionSize)

			binary.Write(f, binary.LittleEndian, decodeMapSize)
			f.Write(decodeMap)

			// With the header in place, write the compressed
			// content.
			f.Write(compressed)
		}
	}

	return nil
}
