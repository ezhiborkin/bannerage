package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

type Storage struct {
	db         *sql.DB
	workerPool *WorkerPool
}

func New(dataSourceName string, workerCount int) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:         db,
		workerPool: NewWorkerPool(workerCount),
	}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgresql.Close"

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

type WorkerPool struct {
	workerCount int
	tasks       chan func() error
}

func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		workerCount: workerCount,
		tasks:       make(chan func() error),
	}

	for i := 0; i < workerCount; i++ {
		go wp.StartWorker(context.Background())
	}

	return wp
}

func (wp *WorkerPool) StartWorker(ctx context.Context) {
	for task := range wp.tasks {
		if err := task(); err != nil {
			log.Printf("Error executing task: %v", err)
		}
	}
}

func (wp *WorkerPool) AddTask(task func() error) {
	wp.tasks <- task
}
