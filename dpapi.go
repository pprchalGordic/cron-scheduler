package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	crypt32                = syscall.NewLazyDLL("crypt32.dll")
	procCryptProtectData   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func dpApiEncrypt(data []byte) ([]byte, error) {
	input := dataBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}
	var output dataBlob

	r, _, err := procCryptProtectData.Call(
		uintptr(unsafe.Pointer(&input)),
		0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(&output)),
	)
	if r == 0 {
		return nil, fmt.Errorf("CryptProtectData failed: %v", err)
	}

	result := make([]byte, output.cbData)
	copy(result, unsafe.Slice(output.pbData, output.cbData))
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(output.pbData)))
	return result, nil
}

func dpApiDecrypt(data []byte) ([]byte, error) {
	input := dataBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}
	var output dataBlob

	r, _, err := procCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&input)),
		0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(&output)),
	)
	if r == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %v", err)
	}

	result := make([]byte, output.cbData)
	copy(result, unsafe.Slice(output.pbData, output.cbData))
	syscall.LocalFree(syscall.Handle(unsafe.Pointer(output.pbData)))
	return result, nil
}

func encryptInteractive() {
	fmt.Print("Enter secret: ")
	var secret string
	fmt.Scanln(&secret)

	encrypted, err := dpApiEncrypt([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "šifrování selhalo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("dpapi:" + base64.StdEncoding.EncodeToString(encrypted))
}

func decryptInteractive() {
	fmt.Print("Zadejte dpapi: ")
	var input string
	fmt.Scanln(&input)

	encoded := strings.TrimPrefix(input, "dpapi:")
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Fprintf(os.Stderr, "není ve fromátu base64: %v\n", err)
		os.Exit(1)
	}

	decrypted, err := dpApiDecrypt(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rozšifrování selhalo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(string(decrypted))
}
