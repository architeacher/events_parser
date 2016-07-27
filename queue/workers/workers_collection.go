package workers

type Collection struct {
	workers          []*Worker
	length           int
}

func NewCollection(workers []*Worker) *Collection {
	return &Collection{
		workers: workers,
		length: len(workers),
	}
}

func (self *Collection) Start() {

	for _, worker := range self.workers {

		go worker.Start()
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