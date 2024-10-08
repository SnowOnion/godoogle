package model

type Resp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

type SearchReq struct {
	Query   string `query:"q"`
	Package string `query:"package"`
}

type SearchResp struct {
	Result []ResultItem `json:"result"`
}

type ResultItem struct {
	ID        int    `json:"id"`
	IDDisplay string `json:"id_display"`
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Pkg       string `json:"pkg"`
	URL       string `json:"url"`
	Signature string `json:"signature"`
	//Distance  int    `json:"distance"` // from query to this result // u.T2 扩展性不如 struct，不太好传出来…… TODO
}
