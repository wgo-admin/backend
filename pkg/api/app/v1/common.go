package v1

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

type IdsRequest []int64

type Pager struct {
	PageNumber int `form:"pageNumber" json:"pageNumber"`
	PageSize   int `form:"pageSize" json:"pageSize"`
}

func (p *Pager) SetDefaultPage() {
	if p.PageNumber == 0 {
		p.PageNumber = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}
}

func (r *Pager) Offset() int {
	r.SetDefaultPage()
	return (r.PageNumber - 1) * r.PageSize
}

func (r *Pager) Limit() int {
	r.SetDefaultPage()
	return r.PageSize
}
