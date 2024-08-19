package main

import (
	"fmt"

	"github.com/superhawk610/bar"
	"github.com/tera-insights/ticrypt-file-copy/copy"
	"github.com/ttacon/chalk"
)

func benchmark(source string, destination string, config *config) error {
	// Benchmark the copy
	progress := make(chan copy.Progress)
	go func() {
		stat := <-progress
		b := bar.NewWithOpts(
			bar.WithDimensions(int(stat.TotalBytes), 20),
			bar.WithFormat(
				fmt.Sprintf(
					" %s copying...%s :percent :bar %s:rate Bytes/s%s :eta",
					chalk.Blue,
					chalk.Reset,
					chalk.Green,
					chalk.Reset,
				),
			),
		)
		for p := range progress {
			b.Update(p.BytesWritten, bar.Context{})
		}
		b.Done()
	}()
	err := copy.NewCopier(source, destination, config.Copy.ChunkSize, progress).Benchmark(copy.Read, copy.Write)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return nil
}