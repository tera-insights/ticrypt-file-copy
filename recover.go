package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/superhawk610/bar"
	"github.com/tera-insights/ticrypt-file-copy/copy"
	"github.com/tera-insights/ticrypt-file-copy/recovery"
	"github.com/ttacon/chalk"
)

const DB = "ticp.db"
const PATH = "/etc/ticp"

func recover() error {
	checkpointer := recovery.NewCheckpointer(filepath.Join(PATH, DB))
	cm := recovery.NewCheckpointManager(checkpointer)
	checkpoints, err := cm.GetInProgressCheckpoints()
	if err != nil {
		return err
	}

	errs := make(chan error, len(checkpoints))
	wg := sync.WaitGroup{}
	wg.Add(len(checkpoints))
	for _, checkpoint := range checkpoints {
		go func(checkpoint *recovery.Checkpoint) {
			defer wg.Done()
			fmt.Printf("Recovering %s\n", checkpoint.SourceFilepath)
			progress := make(chan copy.Progress)
			go func() {
				stat := <-progress
				b := bar.NewWithOpts(
					bar.WithDimensions(int(stat.TotalBytes), 20),
					bar.WithFormat(
						fmt.Sprintf(
							" %scopying... %s :percent :bar %s:rate Bytes/s%s :eta",
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
			err := copy.NewRecoveryCopier(checkpoint.SourceFilepath, checkpoint.DestinationFilePath, checkpoint.ChunkSize, checkpoint.BytesWritten, progress).Copy(copy.Read, copy.Write)
			if err != nil {
				errs <- err
			}

		}(checkpoint)
	}
	wg.Wait()
	close(errs)
	var combinedErr error
	for err := range errs {
		combinedErr = errors.Join(combinedErr, err)
	}
	return combinedErr
}
