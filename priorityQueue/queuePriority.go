// This is the PriorityQueuePackage that further encapsulation by heap
// it is Concurrent security
// You can use priorityQueue.Init() to init a priority queue
// and use (*priorityQueue).PushQueue or (*priorityQueue).PopQueue to operate the heap.Push or heap.Pop

package priorityQueue

import (
	"container/heap"
	"sync"
)

type Entry struct {
	Key      string
	Priority int
}

type PQ []*Entry

type PriorityQueue struct {
	MaxCount int
	HasFull  bool
	Queue    PQ
	Lock     *sync.Mutex
}

func Init(maxCount int) (pq *PriorityQueue) {
	queue := PQ{}
	pq = &PriorityQueue{
		MaxCount: maxCount,
		HasFull:  false,
		Queue:    queue,
		Lock:     &sync.Mutex{},
	}

	heap.Init(pq)
	pq.judgeFull()
	return
}

func (pq *PriorityQueue) judgeFull() bool {
	if pq.Len() < pq.MaxCount {
		pq.HasFull = false
		return false
	} else {
		pq.HasFull = true
		return true
	}
}

// The methods under this line Implements the heap interface

func (pq PriorityQueue) Len() int           { return len(pq.Queue) }
func (pq PriorityQueue) Less(i, j int) bool { return pq.Queue[i].Priority < pq.Queue[j].Priority }
func (pq PriorityQueue) Swap(i, j int)      { pq.Queue[i], pq.Queue[j] = pq.Queue[j], pq.Queue[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	temp := x.(*Entry)
	pq.priorityPush(temp)
}

func (pq *PriorityQueue) Pop() interface{} {
	temp := pq.priorityPop()
	return temp
}

// The methods under this line is further encapsulation

func (pq *PriorityQueue) priorityPush(temp *Entry) {
	pq.Queue = append(pq.Queue, temp)
	pq.judgeFull()
}

func (pq *PriorityQueue) priorityPop() interface{} {
	temp := (pq.Queue)[len(pq.Queue)-1]
	pq.Queue = (pq.Queue)[0 : len(pq.Queue)-1]
	pq.judgeFull()
	return temp
}

// PushQueue is suggest to use in this package
func (pq *PriorityQueue) PushQueue(x interface{}) {
	pq.Lock.Lock()
	for {
		if !pq.HasFull {
			heap.Push(pq, x)
			pq.Lock.Unlock()
			break
		} else {
			pq.Lock.Unlock()
			listenFull := make(chan struct{})
			go func() {
				for {
					if !pq.HasFull {
						listenFull <- struct{}{}
						break
					}
				}
			}()
			// block until queue is not full
			<-listenFull
			pq.Lock.Lock()
			heap.Push(pq, x)
			pq.Lock.Unlock()
			break
		}
	}
}

// PopQueue is suggest to use in this package
func (pq *PriorityQueue) PopQueue() (temp interface{}) {
	pq.Lock.Lock()

	for {
		if pq.Len() != 0 {
			temp = heap.Pop(pq)
			pq.Lock.Unlock()
			return
		} else {
			pq.Lock.Unlock()
			listenEmpty := make(chan struct{})
			go func() {
				for {
					pq.Lock.Lock()
					if pq.Len() != 0 {
						listenEmpty <- struct{}{}
						pq.Lock.Unlock()
						break
					}
					pq.Lock.Unlock()
				}
			}()
			// block until queue is not empty
			<-listenEmpty
			pq.Lock.Lock()
			temp = heap.Pop(pq)
			pq.Lock.Unlock()
			return temp
		}
	}
}
