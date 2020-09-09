package unused

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/PrideSt/otus-golang/hw05_parallel_execution/internal/storage"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidGrtnCnt      = errors.New("invalid goroutin count given")
)

type Task func() error

// Launch runs given function and recover panic if happens.
func (t Task) Launch() (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("task failed with panic %v", p)
		}
	}()

	err = t()

	return err
}

type worker struct {
	id       int
	chTask   chan Task
	chReady  chan<- *worker
	chStop   chan struct{}
	chStoped chan struct{}
	err      error
}

type master struct {
	chReady  chan *worker
	tasks    storage.Interface
	grtnCnt  int
	errLimit int
	errCnt   int
	workers  map[int]*worker
	wmu      sync.Mutex
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(t []Task, grtnCnt int, errLimit int) error {
	if grtnCnt < 1 {
		return fmt.Errorf("%w: expected >1, actual %d", ErrInvalidGrtnCnt, grtnCnt)
	}

	chReady := make(chan *worker, grtnCnt)

	dataCopy := make([]storage.Task, len(t))
	for i, tt := range t {
		dataCopy[i] = storage.Task(tt)
	}

	ts := storage.New(dataCopy)

	m := &master{
		chReady:  chReady,
		tasks:    ts,
		grtnCnt:  grtnCnt,
		errLimit: errLimit,
		workers:  make(map[int]*worker, grtnCnt),
	}

	m.startWorkers()

	return m.runMaster()
}

func (m *master) startWorkers() {
	m.wmu.Lock()
	defer m.wmu.Unlock()

	log.Printf("[master] start all workers\n")
	for i := 0; i < m.grtnCnt; i++ {
		log.Printf("[master] start worker %d [%d/%d]\n", i, i+1, m.grtnCnt)
		m.workers[i] = m.runWorker(i)
	}
}

func (m *master) runWorker(id int) *worker {
	w := &worker{
		id:       id,
		chTask:   make(chan Task),
		chReady:  m.chReady,
		chStop:   make(chan struct{}),
		chStoped: make(chan struct{}),
	}

	go func(w *worker) {
		tN := 0
		log.Printf("[worker-%d] started\n", w.id)
		defer log.Printf("[worker-%d] terminated\n", w.id)
		defer close(w.chStoped)
		defer close(w.chTask)
		for {
			log.Printf("[worker-%d] ready for new tasks...\n", w.id)
			w.chReady <- w
			select {
			case t := <-w.chTask:
				log.Printf("[worker-%d] run new task #%d\n", w.id, tN)
				tN++
				w.err = t.Launch()
				log.Printf("[worker-%d] finish task #%d\n", w.id, tN)
			case <-w.chStop:
				return
			}
		}
	}(w)

	return w
}

func (m *master) stopWorkers() {
	log.Printf("[master] stop all workers\n")
	m.wmu.Lock()
	defer m.wmu.Unlock()

	for _, w := range m.workers {
		w.stop()
	}

	m.workers = nil
}

func (m *master) stopWorker(w *worker) {
	log.Printf("[master] stop worker %d\n", w.id)
	m.wmu.Lock()
	defer m.wmu.Unlock()

	w.stop()
	delete(m.workers, w.id)
}

func (w *worker) stop() {
	close(w.chStop)
	<-w.chStoped
}

func (m *master) runMaster() error {
	defer close(m.chReady)
	for w := range m.chReady {
		// check last error
		if w.err != nil {
			err := w.err
			w.err = nil

			log.Printf("[master] task in worker-%d failed with error: %q", w.id, err)

			m.errCnt++
			if m.errLimit > 0 && m.errCnt == m.errLimit {
				log.Printf("[master] error limit reached [%d/%d]", m.errCnt, m.errLimit)
				m.stopWorkers()

				return ErrErrorsLimitExceeded
			}
			log.Printf("[master] errors [%d/%d]", m.errCnt, m.errLimit)
		}

		// get nex task and push
		task := Task(m.tasks.Next())
		if task == nil {
			log.Printf("[master] there is no tasks for worker %d", w.id)
			m.stopWorker(w)

			m.wmu.Lock()
			if len(m.workers) == 0 {
				return nil
			}
			m.wmu.Unlock()
		} else {
			w.chTask <- task
		}
	}

	return nil
}
