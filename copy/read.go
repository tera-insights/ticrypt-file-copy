package copy

import (
	"fmt"
	"io"
	"os"
)

func Read(copier *Copier) {
	// Read the file
	// fd, err := os.OpenFile(copier.SourceFilepath, os.O_RDONLY|syscall.O_DIRECT, 0644)
	fd, err := os.Open(copier.SourceFilepath)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
	defer fd.Close()
	defer close(copier.ReadDone)

	// Read the file
	for offset := copier.StartingOffset; ; offset += int64(copier.ChunkSize * 1024 * 1024) {
		// Read the file
		n, err := fd.ReadAt(copier.MmapRead, int64(offset))
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
