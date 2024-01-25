package model

type EnrichedMan struct {
	Id         int
	Name       string
	Surname    string
	Patronymic string
	Age        int
	Gender     string
	Country    string
}

func NewEnrichedMan(man *Man, a int, g, c string) *EnrichedMan {
	return &EnrichedMan{
		Name:       man.Name,
		Surname:    man.Surname,
		Patronymic: man.Patronymic,
		Age:        a,
		Gender:     g,
		Country:    c,
	}
}
