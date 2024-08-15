package recovery

type Status string

const (
	// StatusInProgress is the status when the recovery process is in progress
	StatusInProgress Status = "in_progress"
	// StatusComplete is the status when the recovery process is complete
	StatusComplete Status = "complete"
)

// Checkpoint is a struct that represents a checkpoint in the recovery process
type Checkpoint struct {
	SourceFilepath      string `json:"sourceFilepath"`
	DestinationFilePath string `json:"destinationFilePath"`
	ChunkSize           int    `json:"chunkSize"`
	BytesWritten        int64  `json:"bytesWritten"`
	Status              Status `json:"status" gorm:"type:text"`
}

type CheckpointManager struct {
	Checkpointer checkpointer
}

type checkpointer interface {
	CreateCheckpoint(*Checkpoint) error
	GetInProgressCheckpoints() ([]*Checkpoint, error)
}

func NewCheckpointManager(checkpointer checkpointer) *CheckpointManager {
	return &CheckpointManager{
		Checkpointer: checkpointer,
	}
}

func (cm *CheckpointManager) CreateCheckpoint(checkpoint *Checkpoint) error {
	return cm.Checkpointer.CreateCheckpoint(checkpoint)
}

func (cm *CheckpointManager) GetInProgressCheckpoints() ([]*Checkpoint, error) {
	return cm.Checkpointer.GetInProgressCheckpoints()
}
