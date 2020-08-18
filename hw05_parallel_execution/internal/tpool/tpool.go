package tpool

import (
	"container/list"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Task func() error

type TPool struct {
	taskQueue *list.List
	errLimit int
	errCnt int
	workersCnt int
	newTaskPollInterval time.Duration
	chTask chan Task
	chErr chan error
	chTermWorker chan struct{}
	chTermMaster chan struct{}
	chWait chan struct{}
	inProcess int32
	log *log.Logger
	wgMaster *sync.WaitGroup
	wgWorker *sync.WaitGroup
	mu *sync.Mutex
}

func New(sz int, errLimit int, log *log.Logger) *TPool {
	if sz < 1 {
		panic(fmt.Sprintf("given negative TPool size: %d", sz))
	}

	return &TPool{
		taskQueue: list.New(),
		errLimit: errLimit,
		workersCnt: sz,
		newTaskPollInterval: 1000 * time.Millisecond,
		chTask: make(chan Task),
		chErr: make(chan error),
		chTermWorker: make(chan struct{}),
		chTermMaster: make(chan struct{}),
		chWait: nil,
		inProcess: 0,
		log: log,
		wgMaster: &sync.WaitGroup{},
		wgWorker: &sync.WaitGroup{},
		mu: &sync.Mutex{},
	}
}

func (p *TPool) runWorker(workerID int) {
	// p.chTermWorker read only
	// p.chTask read only
	// p.chErr write only, no close
	defer p.wgWorker.Done()
	p.log.Printf("[worker-%d] starts\n", workerID)
	for {
		select{
		case <- p.chTermWorker:
			p.log.Printf("[worker-%d] gracefully termination\n", workerID)
			return
		case t, ok := <- p.chTask:
			if !ok {
				return
			}

			atomic.AddInt32(&p.inProcess, 1)

			err := t(); if err != nil {
				p.log.Printf("[worker-%d] tasks ends with error: %s\n", workerID, err)
				for {
					select {
					case <- p.chTermWorker:
						p.log.Printf("[worker-%d] gracefully termination while error handling\n", workerID)
						return
					case p.chErr <- err:
					}
				}
			} else {
				p.log.Printf("[worker-%d] task done\n", workerID,)
			}

			atomic.AddInt32(&p.inProcess, -1)
		}
	}
}

func (p *TPool) runMaster() {
	// p.chTermMaster read only, no close
	// p.chTask write only, no close
	// p.chErr read and close after closing all workers
	// p.chTermWorker close only, master process close all workers
	defer close(p.chErr) // can close only after all workers terminates, we need waitGroup for it
	defer p.wgMaster.Done()

	p.log.Println("[master] starting")

	if p.chWait != nil {
		p.log.Println("[master] somebody already wait us")
	}

	for i := 1; i <= p.workersCnt; i++ {
		p.log.Printf("[master] launch worker [%d/%d]", i, p.workersCnt)
		p.wgWorker.Add(1)
		go p.runWorker(i)
	}

	var task *Task

	MasterLoop:
	for {
		if task == nil {
			task = p.front()
		}

		// check if any new tasks exists
		if task != nil {
			select{
			case <-p.chTermMaster:
				p.log.Println("[master] gracefully termination")
				break MasterLoop
			case err := <- p.chErr:
				task = nil // skip current task
				if p.handleWorkerError(err) {
					break MasterLoop
				}
			case p.chTask <- *task:
				p.log.Println("[master] send task to worker")
				task = nil
			}
		} else {
			select{
			case <-p.chTermMaster:
				p.log.Println("[master] gracefully termination")
				break MasterLoop
			case err := <- p.chErr:
				if p.handleWorkerError(err) {
					break MasterLoop
				}
			default:
				p.log.Println("[master] task queue is empty, wait new tasks...")
				inProcessNow := atomic.LoadInt32(&p.inProcess)
				if inProcessNow == 0 {
					p.closeWait()
				}
				time.Sleep(p.newTaskPollInterval)
			}
		}
	}

	close(p.chTermWorker)
	p.wgWorker.Wait()
}

func (p *TPool) closeWait() {
	if p.chWait != nil {
		p.log.Println("[master] close wait channel")
		close(p.chWait)
	}
}

func (p *TPool) handleWorkerError(err error) (bool) {
	isTerminate := false
	p.log.Println("[master]", err)
	p.errCnt++

	p.log.Printf("[master] error received: [%d/%d]\n", p.errCnt, p.errLimit)
	if p.errLimit > -1 && p.errCnt >= p.errLimit {
		p.log.Printf("[master] error limit reached, terminate...\n")
		p.closeWait()
		isTerminate = true
	}

	return isTerminate
}

func (p *TPool) Run() {
	p.wgMaster.Add(1)
	go p.runMaster()
}

func (p *TPool) Stop() {
	close(p.chTermMaster)
	p.wgMaster.Wait()
}

// Wait locks add new tasks and wait till all current jobs done.
func (p *TPool) WaitCurrent() {
	// p.mu.Lock()
	// defer p.mu.Unlock()

	p.log.Println("[master] wait till current tasks done...")

	p.chWait = make(chan struct{})

	<- p.chWait
}

func (p *TPool) PushBack(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.taskQueue.PushBack(&t)
}

func (p *TPool) PushFront(t Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.taskQueue.PushFront(&t)
}

func (p *TPool) front() *Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	item := p.taskQueue.Front()
	if item == nil {
		return nil
	}

	p.taskQueue.Remove(item)

	return item.Value.(*Task)
}