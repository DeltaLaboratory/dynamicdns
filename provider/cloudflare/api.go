package cloudflare

import (
	"fmt"
	"time"

	"github.com/imroc/req/v3"

	"github.com/DeltaLaboratory/dynamicdns/dnsapi"
)

type API struct {
	client *req.Client
}

func NewAPI(key string) (*dnsapi.API, error) {
	var iface dnsapi.API
	api := API{
		client: req.NewClient(),
	}
	api.client.SetBaseURL("https://api.cloudflare.com/client/v4/")
	api.client.SetCommonHeader("Authorization", fmt.Sprintf("Bearer %s", key))
	iface = &api
	return &iface, api.verifyToken()
}

func (api *API) verifyToken() error {
	response := new(baseResponse)
	res, err := api.client.R().SetSuccessResult(response).
		SetErrorResult(response).Get("user/tokens/verify")
	if err != nil {
		return fmt.Errorf("request could not be executed: %w", err)
	}
	if res.IsErrorState() {
		base := fmt.Errorf("request failed")
		for _, v := range response.Errors {
			base = fmt.Errorf("%w: [%d:%s]", base, v.Code, v.Message)
		}
		return base
	}
	return nil
}

func (api *API) Create(zoneID string, record dnsapi.Record) error {
	if record.TTL == -1 {
		record.TTL = 1
	}
	response := new(baseResponse)
	res, err := api.client.R().
		SetContentType("application/json").
		SetSuccessResult(response).
		SetErrorResult(response).
		SetBody(createRecordRequest{
			Type:    record.Type,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.TTL,
		}).Post(fmt.Sprintf("zones/%s/dns_records", zoneID))
	if err != nil {
		return fmt.Errorf("request could not be executed: %w", err)
	}
	if res.IsErrorState() {
		base := fmt.Errorf("request failed")
		for _, v := range response.Errors {
			base = fmt.Errorf("%w: [%d:%s]", base, v.Code, v.Message)
		}
		return base
	}
	return nil
}

func (api *API) Update(zoneID string, record dnsapi.Record) error {
	if record.TTL == -1 {
		record.TTL = 1
	}
	id, err := api.findRecord(zoneID, record)
	if err != nil {
		return fmt.Errorf("failed to find record: %w", err)
	}
	response := new(baseResponse)
	res, err := api.client.R().
		SetContentType("application/json").
		SetSuccessResult(response).
		SetErrorResult(response).
		SetBody(patchRecordRequest{
			Type:    record.Type,
			Name:    record.Name,
			Content: record.Content,
			TTL:     record.TTL,
		}).Patch(fmt.Sprintf("zones/%s/dns_records/%s", zoneID, id))
	if err != nil {
		return fmt.Errorf("request could not be executed: %w", err)
	}
	if res.IsErrorState() {
		base := fmt.Errorf("request failed")
		for _, v := range response.Errors {
			base = fmt.Errorf("%w: [%d:%s]", base, v.Code, v.Message)
		}
		return base
	}
	return nil
}

func (api *API) Exists(zoneID string, record dnsapi.Record) bool {
	_, err := api.findRecord(zoneID, record)
	if err != nil {
		return false
	}
	return true
}

func (api *API) findRecord(zoneID string, record dnsapi.Record) (string, error) {
	response := new(listRecordResponse)
	res, err := api.client.R().
		SetSuccessResult(response).
		SetErrorResult(response).
		SetQueryParam("name", record.Name).
		SetQueryParam("type", record.Type).
		Get(fmt.Sprintf("zones/%s/dns_records", zoneID))
	if err != nil {
		return "", fmt.Errorf("request could not be executed: %w", err)
	}
	if res.IsErrorState() {
		base := fmt.Errorf("request failed")
		for _, v := range response.Errors {
			base = fmt.Errorf("%w: [%d:%s]", base, v.Code, v.Message)
		}
		return "", base
	}
	if len(response.Result) == 0 {
		return "", fmt.Errorf("no result matching for filter")
	}
	return response.Result[0].Id, nil
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type createRecordRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type patchRecordRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type baseResponse struct {
	Success  bool       `json:"success"`
	Errors   []apiError `json:"errors"`
	Messages []apiError `json:"messages"`
}

type listRecordResponse struct {
	Success  bool       `json:"success"`
	Errors   []apiError `json:"errors"`
	Messages []apiError `json:"messages"`
	Result   []struct {
		Id         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxiable  bool      `json:"proxiable"`
		Proxied    bool      `json:"proxied"`
		Comment    string    `json:"comment"`
		Tags       []string  `json:"tags"`
		Ttl        int       `json:"ttl"`
		Locked     bool      `json:"locked"`
		ZoneId     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
	} `json:"result"`
}
