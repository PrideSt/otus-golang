package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidGrtnCnt      = errors.New("invalid goroutine count given")
)

type Task func() error

// Launch runs given function and recover panic if happens.
func (t Task) Launch() (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("task failed with panic %v", p)
		}
	}()

	return t()
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, grtnCnt int, errLimit int) error {
	if grtnCnt < 1 {
		return fmt.Errorf("%w: expected >1, actual %d", ErrInvalidGrtnCnt, grtnCnt)
	}

	var errCnt int32
	errLimitInt32 := int32(errLimit)
	chDoneQueue := make(chan struct{}, grtnCnt)
	defer close(chDoneQueue)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for _, tt := range tasks {
		task := tt
		if errLimit < 0 || errLimitInt32 > atomic.LoadInt32(&errCnt) {
			chDoneQueue <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := task.Launch(); err != nil {
					atomic.AddInt32(&errCnt, 1)
				}

				<-chDoneQueue
			}()

			continue
		}

		return ErrErrorsLimitExceeded
	}

	return nil
}
