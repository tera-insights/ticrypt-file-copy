package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func Read(copier *Copier) {
	// Read the file
	fd, err := os.OpenFile(copier.SourceFilepath, os.O_RDONLY|syscall.O_DIRECT, 0644)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
	defer fd.Close()
	defer close(copier.ReadDone)

	// Read the file
	for offset := 0; ; offset += copier.ChunkSize * 1024 {
		// Read the file
		// start := time.Now()
		n, err := fd.ReadAt(copier.MmapRead, int64(offset))
		// t := time.Now()
		// elapsed := t.Sub(start)
		// fmt.Printf("Time taken to read %v\n", elapsed)
		<-copier.WriteDone

		if err == io.EOF {
			copier.MmapRead, copier.MmapWrite = copier.MmapWrite, copier.MmapRead
			copier.ReadDone <- n
			<-copier.WriteDone
			return
		}
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}

		copier.MmapRead, copier.MmapWrite = copier.MmapWrite, copier.MmapRead
		copier.ReadDone <- n
	}
}
