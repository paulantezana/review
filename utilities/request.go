package utilities

import "github.com/paulantezana/review/config"

// ID = cualquier id que se quiera enviar
type Request struct {
	Search      string `json:"search"`
	CurrentPage uint   `json:"current_page"`
	Limit       uint   `json:"limit"`
	ID          uint   `json:"id"`
	ProgramID uint `json:"program_id"`
	SubsidiaryID uint `json:"subsidiary_id"`
}

func (r *Request) Validate() uint {
	con := config.GetConfig()
	if r.Limit == 0 {
		r.Limit = con.Global.Paginate
	}
	if r.CurrentPage == 0 {
		r.CurrentPage = 1
	}
	offset := r.Limit*r.CurrentPage - r.Limit
	return offset
}

// DeleteRequest use in multiple deletes
type DeleteRequest struct {
	Ids []uint `json:"ids"`
}
