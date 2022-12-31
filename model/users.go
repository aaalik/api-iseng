package model

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Dob       string `json:"dob"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type RequestCreateUser struct {
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type RequestUpdateUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type RequestListUser struct {
	Includes *RequestFilterUser `json:"includes"`
	Queries  *Queries           `json:"queries"`
}

type RequestFilterUser struct {
	IDs       []string `json:"ids"`
	Dobs      []string `json:"dobs"`
	CreatedAt *Range   `json:"created_at"`
}

type ResponseUser struct {
	*User
}

type ResponseListUser struct {
	Users []*User        `json:"users"`
	Stats *ResponseStats `json:"stats"`
}

type ResponseStats struct {
	Total int32 `json:"total"`
}
