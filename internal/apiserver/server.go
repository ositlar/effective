package apiserver

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/ositlar/effective/internal/model"
	"github.com/ositlar/effective/internal/sqlstore"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
	store  sqlstore.Store
}

func NewServer(store *sqlstore.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  *store,
	}
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	//ConfigureRouter...
	s.router.HandleFunc("/find", s.handleFind).Methods("GET")
	s.router.HandleFunc("/add", s.handleAdd).Methods("POST")
	s.router.HandleFunc("/delete", s.handleDelete).Methods("DELETE")
	s.router.HandleFunc("/update", s.handleUpdate).Methods("PUT")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) Enrich(m *model.Man) (*model.EnrichedMan, error) {

	type Country struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	}

	type Response struct {
		Countries []Country `json:"country"`
	}

	data := make(map[string]interface{})
	genderURL := &url.URL{
		Scheme:   "https",
		Host:     "api.genderize.io",
		Path:     "/",
		RawQuery: "name=" + m.Name,
	}
	ageURL := &url.URL{
		Scheme:   "https",
		Host:     "api.agify.io",
		Path:     "/",
		RawQuery: "name=" + m.Name,
	}
	countryURL := &url.URL{
		Scheme:   "https",
		Host:     "api.nationalize.io",
		Path:     "/",
		RawQuery: "name=" + m.Name,
	}
	var netClient = http.Client{
		Timeout: time.Second * 10,
	}
	var res *http.Response

	//Age
	//s.logger.Infoln(genderURL.String())
	res, err := netClient.Get(genderURL.String())
	if err != nil {
		s.logger.Infoln("Gender enriching error: ", err)
		//res.Body.Close()
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Infoln("Gender enriching error(reading body): ", err)
		//res.Body.Close()
		return nil, err
	}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		s.logger.Infoln("Gender enriching error(unmarshling body): ", jsonErr)
		//res.Body.Close()
		return nil, jsonErr
	}

	//Age
	res, err = netClient.Get(ageURL.String())
	if err != nil {
		s.logger.Infoln("Age enriching error: ", err)
		//res.Body.Close()
		return nil, err
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		s.logger.Infoln("Age enriching error(reading body): ", err)
		//res.Body.Close()
		return nil, err
	}
	jsonErr = json.Unmarshal(body, &data)
	if jsonErr != nil {
		s.logger.Infoln("Age enriching error(unmarshling body): ", jsonErr)
		//res.Body.Close()
		return nil, jsonErr
	}

	//Country
	var respose Response
	res, err = netClient.Get(countryURL.String())
	if err != nil {
		s.logger.Infoln("Country enriching error: ", err)
		//res.Body.Close()
		return nil, err
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		s.logger.Infoln("Country enriching error(reading body): ", err)
		//res.Body.Close()
		return nil, err
	}
	jsonErr = json.Unmarshal(body, &respose)
	if jsonErr != nil {
		s.logger.Infoln("Country enriching error(unmarshling body): ", jsonErr)
		//res.Body.Close()
		return nil, jsonErr
	}
	//res.Body.Close()
	gender := data["gender"].(string)
	age := int(data["age"].(float64))
	country := respose.Countries[0].CountryID
	em := model.NewEnrichedMan(m, age, gender, country)
	return em, nil
}
