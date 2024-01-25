package apiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ositlar/effective/internal/model"
	"github.com/sirupsen/logrus"
)

func (s *Server) handleFind(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")
	id := r.URL.Query().Get("id")

	switch {
	case name != "":
		m, err := s.store.FindByName("name", name)
		if err != nil {
			s.logger.Log(logrus.InfoLevel, "Wrong get query(name): ", name, err, "\n")
		}
		s.logger.Infoln("Select query success: ", m)
		s.Respond(w, r, http.StatusOK, m)
		return
	case surname != "":
		m, err := s.store.FindByName("surname", surname)
		if err != nil {
			s.logger.Log(logrus.InfoLevel, "Wrong get query(surname): ", surname, err, "\n")
		}
		s.logger.Infoln("Select query success: ", m)
		s.Respond(w, r, http.StatusOK, m)
		return
	case id != "":
		m, err := s.store.FindByName("id", id)
		if err != nil {
			s.logger.Log(logrus.InfoLevel, "Wrong get query(name): ", name, err, "\n")
		}
		s.logger.Infoln("Select query success: ", m)
		s.Respond(w, r, http.StatusOK, m)
		return
	}
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	var m model.Man
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Infoln("Reading body error: ", err)
		s.Respond(w, r, http.StatusBadRequest, "Reading body error")
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		s.logger.Infoln("Unmarshaling body error: ", err)
		s.Respond(w, r, http.StatusInternalServerError, "Unmarshaling body error")
		return
	}
	em, err := s.Enrich(&m)
	if err != nil {
		s.Respond(w, r, http.StatusInternalServerError, err)
		s.logger.Infoln("Enriching error: ", err)
		return
	}
	err = s.store.Insert(em)
	if err != nil {
		s.Respond(w, r, http.StatusInternalServerError, err)
		s.logger.Infoln("Inserting error: ", err)
		return
	}
	s.logger.Infoln("Successful add new EnrichedMan: ", em)
	s.Respond(w, r, http.StatusOK, "Successfully added")
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		s.logger.Infoln("Covertation error: ", err)
		s.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	s.logger.Infoln(id)
	err = s.store.Delete(id)
	if err != nil {
		s.logger.Infoln("Delete error: ", err)
		s.Respond(w, r, http.StatusInternalServerError, err)
		return
	}
	s.logger.Infoln("Delete success: ", id)
	s.Respond(w, r, http.StatusOK, "Delete success: "+fmt.Sprint(id))
}

func (s *Server) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var em model.EnrichedMan
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Infoln("Reading body error: ", err)
		s.Respond(w, r, http.StatusBadRequest, "Reading body error")
		return
	}
	err = json.Unmarshal(body, &em)
	if err != nil {
		s.logger.Infoln("Unmarshaling body error: ", err)
		s.Respond(w, r, http.StatusInternalServerError, "Unmarshaling body error")
		return
	}
	err = s.store.Update(&em)
	if err != nil {
		s.logger.Infoln("Updating error: ", err)
		s.Respond(w, r, http.StatusInternalServerError, "Updating error")
		return
	}
	s.logger.Infoln("Successful update ", em.Id)
	s.Respond(w, r, http.StatusOK, "Successfully updated")

}
