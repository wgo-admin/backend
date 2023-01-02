package v1

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
