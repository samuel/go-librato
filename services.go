package librato

import (
	"fmt"
	"json"
	"os"
)

type Service struct {
	ID       int               `json:"id"`
	Type     string            `json:"type"`
	Settings map[string]string `json:"settings"`
	Title    string            `json:"title"`
}

type ServicesResponse struct {
	Query    QueryResponse `json:"query"`
	Services []Service     `json:"service"`
}

func (s *Service) String() string {
	return fmt.Sprintf("{ID:%d Type:%s Settings:%s Title:%s}", s.ID, s.Type, s.Settings, s.Title)
}

func (r *ServicesResponse) String() string {
	return fmt.Sprintf("{Query:%s Services:%s}", r.Query.String(), r.Services)
}

func (met *Metrics) GetServices() (*ServicesResponse, os.Error) {
	res, err := met.get(librato_metrics_services_api_url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, os.NewError(res.Status)
	}

	var svc ServicesResponse
	jdec := json.NewDecoder(res.Body)
	err = jdec.Decode(&svc)
	return &svc, err
}
