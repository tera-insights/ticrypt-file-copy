package copy

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/sys/unix"
)

const DEFAULT_CHUNK_SIZE = 4 //In MB

type Progress struct {
	BytesWritten int
	TotalBytes   int64
}

type Copier struct {
	CopyID string

	SourceFilepath      string
	DestinationFilePath string
	ChunkSize           int //In MB

	StartingOffset int64

	MmapRead  []byte
	MmapWrite []byte

	ReadDone  chan int
	WriteDone chan int
	Progress  chan Progress
}

func NewCopier(SourceFilepath string, DestinationFilePath string, chunkSize int, progress chan Progress) *Copier {
	if chunkSize <= 0 {
		chunkSize = DEFAULT_CHUNK_SIZE
	}
	copier := &Copier{
		CopyID:              uuid.NewV4().String(),
		SourceFilepath:      SourceFilepath,
		DestinationFilePath: DestinationFilePath,
		ChunkSize:           chunkSize * 1024 * 1024,
		ReadDone:            make(chan int),
		WriteDone:           make(chan int),
		Progress:            progress,
	}
	return copier
}

func NewRecoveryCopier(SourceFilepath string, DestinationFilePath string, chunkSize int, startingOffset int64, progress chan Progress) *Copier {
	if chunkSize <= 0 {
		chunkSize = DEFAULT_CHUNK_SIZE
	}
	copier := &Copier{
		SourceFilepath:      SourceFilepath,
		DestinationFilePath: DestinationFilePath,
		ChunkSize:           chunkSize,
		StartingOffset:      startingOffset,
		ReadDone:            make(chan int),
		WriteDone:           make(chan int),
		Progress:            progress,
	}
	return copier
}

func (c *Copier) Copy(Read func(c *Copier), Write func(c *Copier) <-chan int) error {
	// Log file names
	fmt.Printf("Source: %s\n", c.SourceFilepath)
	fmt.Printf("Destination: %s\n", c.DestinationFilePath)

	mmap, err := unix.Mmap(-1, 0, c.ChunkSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return fmt.Errorf("MMap creation failed %w", err)
	}
	c.MmapRead = mmap

	mmap, err = unix.Mmap(-1, 0, c.ChunkSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return fmt.Errorf("MMap creation failed %w", err)
	}
	c.MmapWrite = mmap

	fmt.Println("MMap Created")

	defer func() {
		err := unix.Munmap(c.MmapRead)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}

		err = unix.Munmap(c.MmapWrite)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			return
		}

		close(c.Progress)
	}()

	fmt.Println("MMap Closed")

	stat, err := os.Stat(c.SourceFilepath)
	if err != nil {
		return fmt.Errorf("Error %v\n", err)
	}

	fmt.Printf("File Size: %v\n", stat.Size())

	progress := Progress{
		BytesWritten: 0,
		TotalBytes:   stat.Size(),
	}
	c.Progress <- progress

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		Read(c)
		fmt.Print("Read Done\n")
	}()

	go func() {
		defer wg.Done()
		writeStats := Write(c)

		for n := range writeStats {
			progress.BytesWritten += n
			c.Progress <- progress
		}
		fmt.Print("Write Done\n")
	}()
	fmt.Print("Waiting for read and write to finish\n")
	wg.Wait()
	fmt.Print("Read and write finished\n")
	return nil
}

func (c *Copier) Benchmark(Read func(c *Copier), Write func(c *Copier) <-chan int) error {
	// Benchmark the copy

	// Benchmark the dd
	fmt.Println("DD Benchmark")
	cmd := exec.Command("dd", []string{"if=/dev/urandom", "of=source", "bs=8M", "count=1280", "iflag=fullblock"}...)
	start := time.Now()
	err := cmd.Run()
	if err != nil {
		return err
	}
	// Because the file is 5GB we multiply by 5*1024 to get the MB
	fmt.Printf("Time taken by dd  %v MB/s \n", 5*1024/time.Since(start).Seconds())

	// Defer the removal of the source file
	defer func() {
		cmd = exec.Command("rm", "source")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// // Benchmark the rsync
	// fmt.Println("rsync Benchmark")
	// cmd = exec.Command("rsync", []string{"source", "destination"}...)
	// start = time.Now()
	// err = cmd.Run()
	// if err != nil {
	// 	return err
	// }
	// // Because the file is 5GB we multiply by 5*1024 to get the MB
	// fmt.Printf("Time taken %v MB/s \n", 5*1024/time.Since(start).Seconds())

	// Defer the removal of the destination file
	defer func() {
		cmd = exec.Command("rm", "destination")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// // Benchmark cp
	// fmt.Println("cp Benchmark")
	// cmd = exec.Command("cp", []string{"source", "destination"}...)
	// start = time.Now()
	// err = cmd.Run()
	// if err != nil {
	// 	return err
	// }
	// // Because the file is 5GB we multiply by 5*1024 to get the MB
	// fmt.Printf("Time taken %v MB/s \n", 5*1024/time.Since(start).Seconds())

	// Benchmark the ticrypt-file-copy
	fmt.Println("ticrypt-file-copy Benchmark")
	err = c.Copy(Read, Write)
	if err != nil {
		return err
	}

	return nil
}
