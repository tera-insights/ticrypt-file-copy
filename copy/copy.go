package copy

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

const DEFAULT_CHUNK_SIZE = 4 //In MB

type Copier struct {
	SourceFilepath      string
	DestinationFilePath string
	MmapRead            []byte
	MmapWrite           []byte
	ChunkSize           int //In MB
	ReadDone            chan int
	WriteDone           chan int
}

func NewCopier(SourceFilepath string, DestinationFilePath string, chunkSize int) *Copier {
	if chunkSize <= 0 {
		chunkSize = DEFAULT_CHUNK_SIZE
	}
	copier := &Copier{
		SourceFilepath:      SourceFilepath,
		DestinationFilePath: DestinationFilePath,
		ChunkSize:           chunkSize,
		ReadDone:            make(chan int),
		WriteDone:           make(chan int),
	}
	return copier
}

func (c *Copier) Copy(Read func(c *Copier), Write func(c *Copier)) error {
	mmap, err := unix.Mmap(0, 0, c.ChunkSize*1024*1024, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return fmt.Errorf("MMap creation failed %w", err)
	}
	c.MmapRead = mmap

	mmap, err = unix.Mmap(0, 0, c.ChunkSize*1024*1024, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return fmt.Errorf("MMap creation failed %w", err)
	}
	c.MmapWrite = mmap

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
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		Read(c)
		fmt.Print("Read Done\n")
	}()
	go func() {
		defer wg.Done()
		Write(c)
		fmt.Print("Write Done\n")
	}()
	wg.Wait()
	return nil
}

func (c *Copier) Benchmark(Read func(c *Copier), Write func(c *Copier)) error {
	// Benchmark the copy

	// Benchmark the dd
	fmt.Println("DD Benchmark")
	cmd := exec.Command("dd", []string{"if=/dev/urandom", "of=source", "bs=64M", "count=16", "iflag=fullblock"}...)
	start := time.Now()
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Time taken by dd to write a 1GB file %v /GB \n", time.Now().Sub(start))

	// Defer the removal of the source file
	defer func() {
		cmd = exec.Command("rm", "source")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// Benchmark the ticrypt-file-copy
	fmt.Println("ticrypt-file-copy Benchmark")
	start = time.Now()
	err = c.Copy(Read, Write)
	if err != nil {
		return err
	}
	fmt.Printf("Time taken %v /GB \n", time.Now().Sub(start))

	// Defer the removal of the destination file
	defer func() {
		cmd = exec.Command("rm", "destination")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// Benchmark the rsync
	fmt.Println("rsync Benchmark")
	cmd = exec.Command("rsync", []string{"source", "destination"}...)
	start = time.Now()
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Time taken %v /GB \n", time.Now().Sub(start))

	return nil
}
