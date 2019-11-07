// Copyright 2019 luomao2000@tom.com. All rights reserved.

// Copy file from possible corrupted file or directory

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	appVersion string = "2019a"
)

var (
	blockSize       int
	retryTimes      int
	dirCount        int
	fileCount       int
	failedFileCount int
	resuedCount     int
	startTime       time.Time
	skipIfFailed    bool
	forceOverwrite  bool
	checkTimestamp  bool
	emptyBuffer     []byte
)

func copyDir(fullName string, f os.FileInfo, outputDir string) {
	fi, err := ioutil.ReadDir(fullName)
	if err != nil {
		println("\ncopyDir error:", err.Error())
	}
	dirName := f.Name()
	for _, f := range fi {
		copy(fullName+string(os.PathSeparator)+f.Name(), f, outputDir+"/"+dirName)
	}
}

func copyFile(fullName string, fi os.FileInfo, outputDir string) {
	t := time.Now()
	fmt.Printf("\r%2d:%02d:%02d %d/%d/%d", t.Hour(), t.Minute(), t.Second(), dirCount, fileCount, failedFileCount)

	outputName := outputDir + string(os.PathSeparator) + fi.Name()
	if !forceOverwrite {
		if st, err := os.Stat(outputName); err == nil && st.Size() == fi.Size() {
			if !checkTimestamp || st.ModTime() == fi.ModTime() {
				resuedCount++
				return
			}
		}
	}

	i, err := os.Open(fullName)
	if err != nil {
		failedFileCount++
		println("\nFailed to open", fullName)
		return
	}
	defer i.Close()

	os.MkdirAll(outputDir, os.ModePerm)
	o, err := os.OpenFile(outputName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode()|0666)
	if err != nil {
		failedFileCount++
		println("\nFailed to create file", err.Error())
		return
	}
	defer o.Close()

	doCopyFile(i, fullName, fi.Size(), o)
}

func doCopyFile(i *os.File, fullName string, fileSize int64, o *os.File) {
	var offset int64
	var failed, needSeek bool
	var err error
	buffer := make([]byte, blockSize)
	for offset = 0; offset < fileSize; offset += int64(blockSize) {
		if needSeek {
			i.Seek(offset, os.SEEK_SET)
			needSeek = false
		}
		size := int(fileSize - offset)
		if size > blockSize {
			size = blockSize
		}

		var bytes int
		for m := 0; m < 1+retryTimes; m++ {
			bytes, err = i.Read(buffer)
			if err == nil {
				break
			}
		}
		if err == nil && bytes > 0 {
			o.Write(buffer[:bytes])
		}
		if err != nil || bytes != size {
			failed = true
			println("\nFailed to read", fullName)
			if skipIfFailed {
				break
			}
			if err != nil || bytes < size {
				if err != nil {
					bytes = 0
				}
				o.Write(emptyBuffer[blockSize-(size-bytes):])
			}
			needSeek = true
		}
	}
	if failed {
		failedFileCount++
	}
}

func copy(fullName string, fi os.FileInfo, outputDir string) {
	if fi == nil {
		println("\nCopy error:", fullName)
	}
	if fi.IsDir() {
		dirCount++
		copyDir(fullName, fi, outputDir)
	} else {
		fileCount++
		copyFile(fullName, fi, outputDir)
	}
	os.Chtimes(outputDir+string(os.PathSeparator)+fi.Name(), fi.ModTime(), fi.ModTime())
}

func main() {
	fmt.Printf("======================================================\n")
	fmt.Printf("badcopy " + appVersion + " is created by luomao2000@tom.com\n")
	fmt.Printf("You can use it to copy corrupted file\n")
	fmt.Printf("======================================================\n")

	var input, output string
	flag.StringVar(&input, "i", "", "*Input directory or file")
	flag.StringVar(&output, "o", "", "*Output directory")
	flag.IntVar(&blockSize, "b", 3, "Block size: 0-512, 1-1K,  2-2K,  3-4K,   4-8K,\n            5-16K, 6-32K, 7-64K, 8-128K, 9-256K\n           ")
	flag.IntVar(&retryTimes, "r", 0, "Retry times: 0 - 9")
	flag.BoolVar(&skipIfFailed, "s", false, "Skip if failed to read file data. Only readed data will be stored.\nLeft data will be skipped.")
	flag.BoolVar(&forceOverwrite, "f", false, "Force overwrite.")
	flag.BoolVar(&checkTimestamp, "c", false, "Check timestamp. File with different timestamp will be treated \nas different file, which will be overwriten.")

	flag.Parse()
	if input == "" {
		println("Please provide input directory or file")
		flag.PrintDefaults()
		return
	}
	if output == "" {
		println("Please provide output directory")
		flag.PrintDefaults()
		return
	}

	if blockSize < 0 || blockSize > 9 {
		println("Block size must between 0 - 9")
		return
	}
	if retryTimes < 0 || retryTimes > 9 {
		println("Retry times must between 0 - 9")
		return
	}
	blockSize = 1 << (9 + blockSize)

	fmt.Printf("         input: %s\n", input)
	fmt.Printf("        output: %s\n", output)
	if blockSize == 512 {
		fmt.Printf("     blockSize: %d\n", blockSize)
	} else {
		fmt.Printf("     blockSize: %dK\n", blockSize/1024)
	}
	fmt.Printf("    retryTimes: %d\n", retryTimes)
	fmt.Printf("  skipIfFailed: %v\n", skipIfFailed)
	fmt.Printf("checkTimestamp: %v\n", checkTimestamp)
	fmt.Printf("forceOverwrite: %v\n", forceOverwrite)
	fmt.Printf("------------------------------------------------------\n")

	fi, err := os.Stat(input)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	startTime = time.Now()
	emptyBuffer = make([]byte, blockSize)
	fmt.Printf("%2d:%02d:%02d started\n", startTime.Hour(), startTime.Minute(), startTime.Second())
	copy(input, fi, output)

	fmt.Printf("\n%6d: directory           \n", dirCount)
	fmt.Printf("%6d: file\n", fileCount)
	fmt.Printf("%6d: copied\n", fileCount-resuedCount)
	fmt.Printf("%6d: reused\n", resuedCount)
	fmt.Printf("%6d: failed\n", failedFileCount)
	duration := int(time.Now().Sub(startTime).Seconds())
	h := duration / (60 * 60)
	m := (duration / 60) % 60
	s := duration % 60
	fmt.Printf("Duration: %2d:%02d:%02d\n", h, m, s)

	fmt.Printf("Done!\n")
}
