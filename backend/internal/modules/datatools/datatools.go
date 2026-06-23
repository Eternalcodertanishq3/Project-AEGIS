package datatools

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// TransformResult holds the result of a data transformation.
type TransformResult struct {
	Input     string `json:"input"`
	Output    string `json:"output"`
	Operation string `json:"operation"`
	Error     string `json:"error,omitempty"`
}

// Transform applies a named operation to the input data.
func Transform(operation, input string) TransformResult {
	result := TransformResult{
		Input:     input,
		Operation: operation,
	}

	switch strings.ToLower(operation) {
	// Encoding
	case "base64-encode":
		result.Output = base64.StdEncoding.EncodeToString([]byte(input))
	case "base64-decode":
		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			result.Error = fmt.Sprintf("invalid base64: %v", err)
		} else {
			result.Output = string(decoded)
		}
	case "url-encode":
		result.Output = url.QueryEscape(input)
	case "url-decode":
		decoded, err := url.QueryUnescape(input)
		if err != nil {
			result.Error = fmt.Sprintf("invalid url encoding: %v", err)
		} else {
			result.Output = decoded
		}
	case "hex-encode":
		result.Output = hex.EncodeToString([]byte(input))
	case "hex-decode":
		decoded, err := hex.DecodeString(input)
		if err != nil {
			result.Error = fmt.Sprintf("invalid hex: %v", err)
		} else {
			result.Output = string(decoded)
		}

	// Hashing
	case "md5":
		h := md5.Sum([]byte(input))
		result.Output = hex.EncodeToString(h[:])
	case "sha1":
		h := sha1.Sum([]byte(input))
		result.Output = hex.EncodeToString(h[:])
	case "sha256":
		h := sha256.Sum256([]byte(input))
		result.Output = hex.EncodeToString(h[:])
	case "sha512":
		h := sha512.Sum512([]byte(input))
		result.Output = hex.EncodeToString(h[:])

	// Text transforms
	case "uppercase":
		result.Output = strings.ToUpper(input)
	case "lowercase":
		result.Output = strings.ToLower(input)
	case "reverse":
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		result.Output = string(runes)
	case "char-count":
		result.Output = fmt.Sprintf("Characters: %d, Words: %d, Lines: %d",
			len([]rune(input)),
			len(strings.Fields(input)),
			strings.Count(input, "\n")+1,
		)
	case "rot13":
		result.Output = strings.Map(func(r rune) rune {
			switch {
			case r >= 'A' && r <= 'Z':
				return 'A' + (r-'A'+13)%26
			case r >= 'a' && r <= 'z':
				return 'a' + (r-'a'+13)%26
			}
			return r
		}, input)

	default:
		result.Error = fmt.Sprintf("unknown operation: %s", operation)
	}

	return result
}

// ListOperations returns all available operations grouped by category.
func ListOperations() map[string][]OperationInfo {
	return map[string][]OperationInfo{
		"Encoding": {
			{ID: "base64-encode", Name: "Base64 Encode", Description: "Encode text to Base64"},
			{ID: "base64-decode", Name: "Base64 Decode", Description: "Decode Base64 to text"},
			{ID: "url-encode", Name: "URL Encode", Description: "Percent-encode text for URLs"},
			{ID: "url-decode", Name: "URL Decode", Description: "Decode percent-encoded text"},
			{ID: "hex-encode", Name: "Hex Encode", Description: "Encode text to hexadecimal"},
			{ID: "hex-decode", Name: "Hex Decode", Description: "Decode hexadecimal to text"},
		},
		"Hashing": {
			{ID: "md5", Name: "MD5", Description: "Compute MD5 hash"},
			{ID: "sha1", Name: "SHA-1", Description: "Compute SHA-1 hash"},
			{ID: "sha256", Name: "SHA-256", Description: "Compute SHA-256 hash"},
			{ID: "sha512", Name: "SHA-512", Description: "Compute SHA-512 hash"},
		},
		"Text": {
			{ID: "uppercase", Name: "Uppercase", Description: "Convert to uppercase"},
			{ID: "lowercase", Name: "Lowercase", Description: "Convert to lowercase"},
			{ID: "reverse", Name: "Reverse", Description: "Reverse text"},
			{ID: "char-count", Name: "Character Count", Description: "Count characters, words, lines"},
			{ID: "rot13", Name: "ROT13", Description: "Apply ROT13 cipher"},
		},
	}
}

// OperationInfo describes an available operation.
type OperationInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
