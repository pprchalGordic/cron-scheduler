package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func LogRotateStart(srcRoot, dstRoot string) {
	srcInfo, err := os.Stat(srcRoot)
	if err != nil || !srcInfo.IsDir() {
		return
	}

	dstInfo, err := os.Stat(dstRoot)
	if err != nil || !dstInfo.IsDir() {
		return
	}

	entries, err := os.ReadDir(srcRoot)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		srcDir := filepath.Join(srcRoot, entry.Name())
		if !canBeProcessed(srcDir) {
			continue
		}

		dstZipFile := filepath.Join(dstRoot, srcInfo.Name()+".zip")

		if _, err := os.Stat(dstZipFile); err == nil {
			appendToZip(dstZipFile, srcDir)
		} else {
			createZipFromDir(dstZipFile, srcDir)
		}

		// delete files and directory
		files, _ := os.ReadDir(srcDir)
		for _, f := range files {
			os.Remove(filepath.Join(srcDir, f.Name()))
		}
		os.RemoveAll(srcDir)
	}
}

func canBeProcessed(srcDir string) bool {
	name := filepath.Base(srcDir)
	t := time.Now().AddDate(0, 0, -1)

	// split by _ or -
	var parts []string
	if strings.Contains(name, "_") {
		parts = strings.Split(name, "_")
	} else {
		parts = strings.Split(name, "-")
	}

	if len(parts) != 3 {
		return false
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return false
	}

	dirDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return dirDate.Before(t)
}

func appendToZip(zipPath, srcDir string) {
	zipFile, err := os.OpenFile(zipPath, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer zipFile.Close()

	info, _ := zipFile.Stat()
	reader, err := zip.NewReader(zipFile, info.Size())
	if err != nil {
		return
	}

	// Create temp file with updated archive
	tmpPath := zipPath + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return
	}

	writer := zip.NewWriter(tmpFile)

	// Copy existing entries
	for _, item := range reader.File {
		rc, err := item.Open()
		if err != nil {
			continue
		}
		w, err := writer.Create(item.Name)
		if err != nil {
			rc.Close()
			continue
		}
		io.Copy(w, rc)
		rc.Close()
	}

	// Add new files
	files, _ := os.ReadDir(srcDir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		entryName := fmt.Sprintf("%s_%d", f.Name(), time.Now().UnixNano())
		w, err := writer.Create(entryName)
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, f.Name()))
		if err != nil {
			continue
		}
		w.Write(data)
	}

	writer.Close()
	tmpFile.Close()
	zipFile.Close()

	os.Rename(tmpPath, zipPath)
}

func createZipFromDir(zipPath, srcDir string) {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	files, _ := os.ReadDir(srcDir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		w, err := writer.Create(f.Name())
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, f.Name()))
		if err != nil {
			continue
		}
		w.Write(data)
	}
}
