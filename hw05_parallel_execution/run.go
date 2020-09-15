package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"log"
	"sync"
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

func worker(id int, chTasks <-chan Task, chErrors chan<- error, wg *sync.WaitGroup) {
	log.Printf("[%d] start worker", id)
	defer wg.Done()
	for t := range chTasks {
		log.Printf("[%d] task received", id)
		if err := t.Launch(); err != nil {
			log.Printf("[%d] error happened", id)
			chErrors <- err
		}
	}
	log.Printf("[%d] task chan closed, terminate worker", id)
}

func handleError(errCnt, errLimit int) error {
	log.Printf("[main] err received %d/%d", errCnt, errLimit)
	if errLimit >= 0 && errCnt >= errLimit {
		log.Printf("[main] err limit reached, terminate")
		return ErrErrorsLimitExceeded
	}

	return nil
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, grtnCnt int, errLimit int) error {
	if grtnCnt < 1 {
		return fmt.Errorf("%w: expected > 0, actual %d", ErrInvalidGrtnCnt, grtnCnt)
	}

	var errCnt int

	// we must use buffered channel, otherwise while we waiting running goroutines termination
	// nobody reads from chErrors and running goroutines will blocked on writing to them
	// and never ends
	chErrors := make(chan error, grtnCnt)
	defer close(chErrors)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	chTasks := make(chan Task)
	defer close(chTasks)

	for i := 0; i < grtnCnt; i++ {
		wg.Add(1)
		go worker(i, chTasks, chErrors, wg)
	}

	i := 0
	for i < len(tasks) {
		select {
		case <-chErrors:
			errCnt++
			if err := handleError(errCnt, errLimit); err != nil {
				return err
			}
		default:
		}

		select {
		case <-chErrors:
			errCnt++
			if err := handleError(errCnt, errLimit); err != nil {
				return err
			}
		case chTasks <- tasks[i]:
			i++
			log.Printf("[main] task %d pushed", i)
		}
	}

	return nil
}
