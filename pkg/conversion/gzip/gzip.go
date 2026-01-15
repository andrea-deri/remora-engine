package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"

	"remora/pkg/customerror"
)

// CompressedStringToPlainString decodes a base64-encoded gzip-compressed string
// and returns the decompressed plain string.
//
// The function performs the following steps:
//  1. Decodes the input string from base64 to get gzip-compressed bytes.
//  2. Creates a gzip.Reader to decompress the byte stream.
//  3. Reads all decompressed data into a byte slice.
//  4. Converts the byte slice to a string and returns it.
//
// Returns a custom error if any of the decoding or decompression steps fail.
func CompressedStringToPlainString(sourceString string) (string, error) {

	gzippedBytes, conversionError := base64.StdEncoding.DecodeString(sourceString)
	if conversionError != nil {
		return "", customerror.NewError(customerror.ErrorGZipNotDecodable, conversionError)
	}

	byteReader := bytes.NewReader(gzippedBytes)
	gzipReader, conversionError := gzip.NewReader(byteReader)
	if conversionError != nil {
		return "", customerror.NewError(customerror.ErrorGZipNotDecodable, conversionError)
	}

	byteArrayResult, conversionError := io.ReadAll(gzipReader)
	if conversionError != nil {
		return "", customerror.NewError(customerror.ErrorGZipNotDecodable, conversionError)
	}

	return string(byteArrayResult), nil
}
