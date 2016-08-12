package workers

import "sync"

type Collection struct {
	workers      []*Worker
	length       int
	isGrouped    bool
}

func NewCollection(workers []*Worker, isGrouped bool) *Collection {
	return &Collection{
		workers: workers,
		length: len(workers),
		isGrouped: isGrouped,
	}
}

func (self *Collection) Start() {

	var wg sync.WaitGroup

	if self.isGrouped {
		wg = sync.WaitGroup{}
	}

	for _, worker := range self.workers {

		if self.isGrouped {
			wg.Add(1)
		}

		go worker.Start(wg)
	}

	if self.isGrouped {
		wg.Wait()
	}
}

func (self *Collection) Pause() {

	for _, worker := range self.workers {

		worker.Pause()
	}
}

func (self *Collection) Stop() {

	for _, worker := range self.workers {

		worker.Stop()
	}
}

func (self *Collection) Restart() {

	self.Stop()
	self.Start()
}