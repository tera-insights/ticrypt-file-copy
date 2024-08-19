package recovery

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Checkpointer struct {
	dbPath string
	db     *gorm.DB
}

func NewCheckpointer(dbPath string) checkpointer {
	return &Checkpointer{
		dbPath: dbPath,
	}
}

func (c *Checkpointer) getDB() (*gorm.DB, error) {
	if c.db == nil {
		_, err := os.Stat(c.dbPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			err = os.MkdirAll(filepath.Dir(c.dbPath), 0755)
			if err != nil {
				return nil, err
			}
		}

		db, err := gorm.Open(sqlite.Open(c.dbPath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		c.db = db
	}

	return c.db, nil
}

func (c *Checkpointer) CreateCheckpoint(checkpoint *Checkpoint) error {
	db, err := c.getDB()
	if err != nil {
		return err
	}

	err = db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(checkpoint).Error
	return c.db.Create(checkpoint).Error
}

func (c *Checkpointer) GetCheckpoints() ([]*Checkpoint, error) {
	db, err := c.getDB()
	if err != nil {
		return nil, err
	}

	var checkpoint []*Checkpoint
	err = db.Find(&checkpoint).Error
	return checkpoint, err
}

func (c *Checkpointer) GetInProgressCheckpoints() ([]*Checkpoint, error) {
	db, err := c.getDB()
	if err != nil {
		return nil, err
	}

	var checkpoint []*Checkpoint
	err = db.Where("status = ?", StatusInProgress).First(&checkpoint).Error
	return checkpoint, err
}
