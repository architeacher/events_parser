package jobs
// Todo: Use one collection instead of different collections.
type Collection struct {
	jobs []Job
	length int
}

func NewCollection(jobs []Job) *Collection {
	return &Collection{
		jobs: jobs,
		length: len(jobs),
	}
}

func (self *Collection) GetLength() int {
	return self.length
}
