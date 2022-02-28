package job

func New() *Job {
	return &Job{}
}

type Job struct {
	Type     string
	Key      string
	Name     string
	Children []*Job
}

type Step func() (Step, error)

func (j *Job) Get() Step {
	return func() (Step, error) {
		return nil, nil
	}
}
