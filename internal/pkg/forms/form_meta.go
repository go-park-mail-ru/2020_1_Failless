package forms

type MetaForm struct {
	Uid    int      `json:"uid"`
	Tags   []int    `json:"tags"`
	About  string   `json:"about"`
	Social []string `json:"social"`
	Photos []EImage `json:"photos"`
}
