package jobs

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

