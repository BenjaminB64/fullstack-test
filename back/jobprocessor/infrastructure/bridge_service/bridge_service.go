package bridge_service

import (
	"encoding/json"
	"errors"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/domain"
	"github.com/BenjaminB64/fullstack-test/back/jobprocessor/infrastructure/logger"
	"net/http"
	"net/url"
	"time"
)

type BridgeService struct {
	httpClient *http.Client
	logger     *logger.Logger
}

func NewBridgeService(logger *logger.Logger) domain.BridgeService {
	httpClient := &http.Client{}
	httpClient.Timeout = 10 * time.Second

	return &BridgeService{
		httpClient: httpClient,
		logger:     logger,
	}
}

func (bs *BridgeService) GetBridgeSchedule() (*domain.BridgeSchedule, error) {
	// get from bordeaux metropole open data
	getUrl := url.URL{}
	getUrl.Scheme = "https"
	getUrl.Host = "opendata.bordeaux-metropole.fr"
	getUrl.Path = "/api/explore/v2.1/catalog/datasets/previsions_pont_chaban/records"

	// https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/api/
	// add query params
	query := url.Values{}
	query.Set("where", "date_passage > now() and date_passage < now(days=15)")
	query.Set("order_by", "date_passage")
	query.Set("limit", "-1")
	query.Set("timezone", "Europe/Paris")
	getUrl.RawQuery = query.Encode()

	get, err := bs.httpClient.Get(getUrl.String())
	if err != nil {
		return nil, err
	}

	defer get.Body.Close()

	// parse response
	if get.StatusCode != http.StatusOK {
		bs.logger.Error("error getting bridge", "status", get.StatusCode)
		return nil, errors.New("error getting bridge")
	}

	var response OpenDataBordeauxResponse
	err = json.NewDecoder(get.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	bridge := &domain.BridgeSchedule{
		Closures: make([]domain.BridgeClosure, 0, len(response.Results)),
	}

	for _, result := range response.Results {
		if result.ClosureTime != "" && result.ReopeningTime != "" {
			closure, reopen, err := BridgeOpenCloseTimes(result.PassageDate, result.ClosureTime, result.ReopeningTime)
			if err != nil {
				return nil, err
			}
			bridge.Closures = append(bridge.Closures, domain.BridgeClosure{
				BoatName:   result.BoatName,
				CloseTime:  closure,
				ReopenTime: reopen,
			})
		}
	}

	return bridge, nil

}

// BridgeOpenCloseTimes parses the passage date, closure and reopen time and returns the closure and reopen datetimes
// (take into account the date of the passage, and the fact that the bridge can close at 23:59 and reopen at 00:01)
func BridgeOpenCloseTimes(passageDate, closureTime, reopenTime string) (closure time.Time, reopen time.Time, err error) {
	frLocale, _ := time.LoadLocation("Europe/Paris")

	closure, err = time.Parse("15:04", closureTime)
	if err != nil {
		err = errors.Join(errors.New("error parsing close time"), err)
		return
	}

	reopen, err = time.Parse("15:04", reopenTime)
	if err != nil {
		err = errors.Join(errors.New("error parsing open time"), err)
		return
	}

	var passage time.Time
	passage, err = time.Parse("2006-01-02", passageDate)
	if err != nil {
		err = errors.Join(errors.New("error parsing passage date"), err)
		return
	}

	reopen = time.Date(passage.Year(), passage.Month(), passage.Day(), reopen.Hour(), reopen.Minute(), reopen.Second(), reopen.Nanosecond(), frLocale)
	closure = time.Date(passage.Year(), passage.Month(), passage.Day(), closure.Hour(), closure.Minute(), closure.Second(), closure.Nanosecond(), frLocale)

	if closure.After(reopen) {
		reopen.Add(24 * time.Hour)
	}

	return
}

type OpenDataBordeauxResponse struct {
	TotalCount int `json:"total_count"`
	Results    []struct {
		BoatName         string `json:"bateau"`
		PassageDate      string `json:"date_passage"`
		ClosureTime      string `json:"fermeture_a_la_circulation"`
		ReopeningTime    string `json:"re_ouverture_a_la_circulation"`
		ClosureType      string `json:"type_de_fermeture"`
		TotalClosure     string `json:"fermeture_totale"`
		ClosureDirection string `json:"sens_fermeture"`
	} `json:"results"`
}
