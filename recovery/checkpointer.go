package recovery

import (
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
