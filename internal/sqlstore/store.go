package sqlstore

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq" //...
	"github.com/ositlar/effective/internal/model"
	"github.com/sirupsen/logrus"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Insert(u *model.EnrichedMan) error {
	logrus.Infoln("Insert into people_info: ", u)
	if err := s.db.QueryRow(
		"INSERT INTO people_info (name, surname, pathronymic, age, sex, nation) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		u.Name, u.Surname, u.Patronymic, u.Age, u.Gender, u.Country).Scan(&u.Id); err != nil {
		return errors.New("Insert error: " + err.Error())
	}
	return nil
}

func (s *Store) FindByName(key, value string) ([]*model.EnrichedMan, error) {
	res := make([]*model.EnrichedMan, 0)
	rows, err := s.db.Query("SELECT id, name, surname, pathronymic, age, sex, nation FROM people_info WHERE $1 = $2", key, value)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			name, surname, pathronymic, gender, nation string
			id, age                                    int
		)
		err = rows.Scan(&id, &name, &surname, &pathronymic, &age, &gender, &nation)
		if err != nil {
			return nil, err
		}
		m := &model.EnrichedMan{
			Id:         id,
			Name:       name,
			Surname:    surname,
			Patronymic: pathronymic,
			Gender:     gender,
			Country:    nation,
			Age:        age,
		}
		res = append(res, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Store) Delete(id int) error {

	_, err := s.db.Query("DELETE FROM people_info WHERE id = $1", id)
	if err != nil {
		return errors.New("Delete error: " + err.Error())
	}
	return nil

}

func (s *Store) Update(u *model.EnrichedMan) error {
	if err := s.Delete(u.Id); err != nil {
		return err
	}
	if err := s.Insert(u); err != nil {
		return err
	}

	return nil
}
