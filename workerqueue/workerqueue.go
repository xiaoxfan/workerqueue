/*
@Time : 2020/4/9 3:56 PM
*/
package workerqueue

import (
	"fmt"
)

var defaultWorkQueueSize = 100
var defaultWorkers = 2

type WorkerQueue struct {
	WorkerQueue chan chan Work // worker的work channel
	WorkQueue   chan Work
	nWorker     int
}

func New(nWorkers int) *WorkerQueue {
	if nWorkers < 1 {
		nWorkers = defaultWorkers
	}
	return &WorkerQueue{WorkerQueue: make(chan chan Work), nWorker: nWorkers, WorkQueue: make(chan Work, defaultWorkQueueSize)}
}

type Work func()

// 添加新任务到任务队列
func (q *WorkerQueue) SubmitWork(work Work) {
	q.WorkQueue <- work
}

type Worker struct {
	ID          int
	Work        chan Work
	WorkerQueue chan chan Work
	QuitChan    chan bool
}

func NewWorker(ID int, workerQueue chan chan Work) *Worker {
	return &Worker{
		ID:          ID,
		Work:        make(chan Work),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
}

// 把worker负责接收work的channel放到WorkerQueue
// WorkerQueue的Start方法会从WorkQueue中取Work，然后从WorkerQueue中取worker负责接收work的channel
// 然后把work放入该worker的work channel
// worker从work channel中取出work 然后开始执行work
// 重复上述过程 每个worker同时只执行一个任务，所以Work channel不需要缓冲
func (w *Worker) Start() {
	go func() {
		for {
			// 把该worker接收work的channel放到WorkerQueue，等待WorkerQueue.Start方法往channel中发送任务
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				fmt.Printf("worker%d: working\n", w.ID)
				work()
			case <-w.QuitChan:
				fmt.Printf("worker%d: stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func (q *WorkerQueue) Start() {
	// 启动
	for i := 0; i < q.nWorker; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, q.WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-q.WorkQueue:
				//fmt.Println("Received work")
				go func() {
					worker := <-q.WorkerQueue
					//fmt.Println("Dispatching work")
					worker <- work
				}()
			}
		}
	}()
}