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

const cryptProtectLocalMachine = 0x4

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
		0, 0, 0, 0, cryptProtectLocalMachine,
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
		0, 0, 0, 0, cryptProtectLocalMachine,
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

func readPassword() string {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	readConsoleInput := kernel32.NewProc("ReadConsoleInputW")

	handle, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	var oldMode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&oldMode)))
	setConsoleMode.Call(uintptr(handle), 0) // disable all input processing

	type inputRecord struct {
		eventType uint16
		_         uint16
		keyDown   int32
		repeat    uint16
		vkCode    uint16
		scanCode  uint16
		char      uint16
		state     uint32
	}

	var result []byte
	for {
		var rec inputRecord
		var read uint32
		readConsoleInput.Call(uintptr(handle), uintptr(unsafe.Pointer(&rec)), 1, uintptr(unsafe.Pointer(&read)))
		if rec.eventType != 1 || rec.keyDown == 0 {
			continue
		}
		ch := rec.char
		if ch == 13 { // Enter
			break
		}
		if ch == 8 { // Backspace
			if len(result) > 0 {
				result = result[:len(result)-1]
				fmt.Print("\b \b")
			}
			continue
		}
		if ch >= 32 {
			result = append(result, byte(ch))
			fmt.Print("*")
		}
	}
	fmt.Println()

	setConsoleMode.Call(uintptr(handle), uintptr(oldMode))
	return string(result)
}

func encryptInteractive() {
	fmt.Print("Zadejte heslo: ")
	secret := readPassword()

	encrypted, err := dpApiEncrypt([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "šifrování selhalo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("dpapi:" + base64.StdEncoding.EncodeToString(encrypted))
}

func decryptInteractive(path string) {
	var input string
	if path != "" {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nelze načíst soubor %s: %v\n", path, err)
			os.Exit(1)
		}
		input = strings.TrimSpace(string(content))
	} else {
		fmt.Print("Zadejte dpapi: ")
		fmt.Scanln(&input)
	}

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
