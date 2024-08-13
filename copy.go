package main

import (
	"fmt"
	"sync"

	"golang.org/x/sys/unix"
)

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
	copier := &Copier{
		SourceFilepath:      SourceFilepath,
		DestinationFilePath: DestinationFilePath,
		ChunkSize:           chunkSize,
		ReadDone:            make(chan int),
		WriteDone:           make(chan int),
	}
	return copier
}

func (c *Copier) Copy() error {
	mmap, err := unix.Mmap(0, 0, c.ChunkSize*1024, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		return fmt.Errorf("MMap creation failed %w", err)
	}
	c.MmapRead = mmap

	mmap, err = unix.Mmap(0, 0, c.ChunkSize*1024, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_PRIVATE|unix.MAP_ANON)
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
