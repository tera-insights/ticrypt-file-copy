package main

import (
	"fmt"
	"os"
	"syscall"
)

func Write(copier *Copier) {
	// Write the file
	fd, err := os.OpenFile(copier.DestinationFilePath, os.O_WRONLY|os.O_CREATE|syscall.O_DIRECT, 0644)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
	defer fd.Close()
	defer close(copier.WriteDone)
	// Write the file
	copier.WriteDone <- 0
	for offset := 0; ; offset += copier.ChunkSize * 1024 {
		// Write the file
		if n, ok := <-copier.ReadDone; ok {
			// start := time.Now()
			_, err := fd.WriteAt(copier.MmapWrite[:n], int64(offset))
			// t := time.Now()
			// elapsed := t.Sub(start)
			// fmt.Printf("Time taken to write %v\n", elapsed)
			if err != nil {
				fmt.Printf("Error %v\n", err)
				return
			}
			copier.WriteDone <- n
		} else {
			return
		}
	}
}
