package model

type Man struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func NewMan(n, s, p string) *Man {
	return &Man{
		Name:       n,
		Surname:    s,
		Patronymic: p,
	}
}
