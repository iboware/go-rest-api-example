package model

type MDBRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type Record struct {
	Key        string `json:"key"`
	CreatedAt  string `json:"createdAt"`
	TotalCount int    `json:"totalCount"`
}

type MDBResponse struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Records []Record `json:"records"`
}
