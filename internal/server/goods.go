package server

import (
	"encoding/json"
	"github.com/0ndreu/test_kamaz/internal/models"
	"io/ioutil"
	"net/http"
)

func (s *server) CreateMultipleGoods(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if body != nil {
		defer r.Body.Close()
	}
	if err != nil {
		s.log.Error().Msgf("CreateMultipleGoods employee body not valid error:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var goods []models.Good
	err = json.Unmarshal(body, &goods)
	if err != nil {
		s.log.Error().Msgf("CreateMultipleGoods unmarshal body error:%s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.svc.CreateMultipleGoods(s.ctx, goods)
}
