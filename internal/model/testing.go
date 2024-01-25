package model

func Test() *EnrichedMan {
	return &EnrichedMan{
		Name:       "Johnsdfsdf",
		Surname:    "Smith",
		Age:        40,
		Patronymic: "Johnson",
		Gender:     "Male",
		Country:    "US",
	}
}
