package api

type Res struct {
	Status   int         `json:"status"`
	Response interface{} `json:"res"`
}
