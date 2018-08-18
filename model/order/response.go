package order

type RespCreateTransaction struct {
	ID int `json:"id"`
	ChangeBalance int `json:"change_balance"`
	Desc string `json:"description"`
	Finished *bool `json:"finished"`
	User string `json:"user"`
}