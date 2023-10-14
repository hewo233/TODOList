package model

// 标准化响应
type Resp struct {
	Code int  `json:"code"`
	Data TODO `json:"data"`
}

// todolist的结构体
type TODO struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Done     bool   `json:"done"`
	Exist    bool   `json:"exist"`
	Email    string `json:"email"` //对应每个客户
	Tag      string `json:"tag"`
	Deadline string `json:"deadline"`
}
