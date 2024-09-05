package copy

import (
	"fmt"
	"os"
)

func Write(copier *Copier) <-chan int {
	// Write the file
	stats := make(chan int)
	go func() {
		defer close(stats)
		// fd, err := os.OpenFile(copier.DestinationFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC|syscall.O_DIRECT, 0666)
		fd, err := os.Create(copier.DestinationFilePath)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}
		defer fd.Close()
		defer close(copier.WriteDone)
		// Write the file
		copier.WriteDone <- 0
		for offset := copier.StartingOffset; ; offset += int64(copier.ChunkSize) {
			// Write the file
			if n, ok := <-copier.ReadDone; ok {
				_, err := fd.WriteAt(copier.MmapWrite[:n], int64(offset))
				if err != nil {
					fmt.Printf("Error %v\n", err)
					return
				}
				select {
				case stats <- n:
				default:
				}
				copier.WriteDone <- n
			} else {
				return
			}
		}
	}()
	return stats
}
