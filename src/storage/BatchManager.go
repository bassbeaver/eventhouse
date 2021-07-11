package storage

import (
	"database/sql"
	"time"
)

type batchManager struct {
	maxEntitiesInBatch int
	maxBatchLifeTime   time.Duration
	dbConnect          *sql.DB
	appendChan         chan *batchEntity
}

func (bs *batchManager) startListener() {
	activeBatch := bs.createNewBatch()

	go func() {
		for <-activeBatch.closedChan {
			activeBatch = bs.createNewBatch()
		}
	}()
}

func (bs *batchManager) createNewBatch() *batch {
	return newBatch(bs.maxEntitiesInBatch, bs.maxBatchLifeTime, bs.dbConnect, bs.appendChan)
}

func (bs *batchManager) Append(batchEntityObj *batchEntity) {
	bs.appendChan <- batchEntityObj
}

func newBatchManager(
	maxEntitiesInBatch int,
	maxBatchLifeTime time.Duration,
	dbConnect *sql.DB,
) *batchManager {
	bs := &batchManager{
		maxEntitiesInBatch: maxEntitiesInBatch,
		maxBatchLifeTime:   maxBatchLifeTime,
		dbConnect:          dbConnect,
		appendChan:         make(chan *batchEntity),
	}
	bs.startListener()

	return bs
}
