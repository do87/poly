package plan

type Plan struct {
	Type string
	Key  string
	Name string
}

type Job func() (Job, error)

func New() *Plan {
	return &Plan{}
}

func (p *Plan) Start() Job {
	return func() (Job, error) {
		return nil, nil
	}
}
