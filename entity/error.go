package entity

type ErrResponse struct {
	SourceFunction string `json:"source_function"`
	Detail         string `json:"err_detail"`
}
