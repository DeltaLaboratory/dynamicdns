package application

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DeltaLaboratory/go-ipify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/DeltaLaboratory/dynamicdns/config"
	"github.com/DeltaLaboratory/dynamicdns/dnsapi"
	"github.com/DeltaLaboratory/dynamicdns/provider"
)

type DDNSService struct {
}

func (d *DDNSService) Run() {
	prov := provider.NewProvider()
	logger := log.With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load config")
		return
	}
	var sid int
	for _, v := range cfg.DDNS.Service {
		sid++
		s := service{
			client:   nil,
			interval: time.Duration(cfg.Interval) * time.Minute,
			config:   v,
			logger:   logger.With().Str("service", fmt.Sprintf("%s-%d", v.Provider, sid)).Logger(),
		}
		if err := s.Init(prov); err != nil {
			s.logger.Error().Err(err).Msg("failed to initialize service")
			return
		} else {
			s.logger.Info().Msg("initialized service")
		}
		go s.Start()
	}
	select {}
}

type service struct {
	client   *dnsapi.API
	interval time.Duration
	config   config.Service
	logger   zerolog.Logger
}

func (s *service) Init(provider *provider.Provider) error {
	creator, err := provider.GetProvider(s.config.Provider)
	if err != nil {
		return fmt.Errorf("failed to intialize service: not supported provider: %w", err)
	}
	s.client, err = creator(s.config.APIKey)
	if err != nil {
		return fmt.Errorf("failed to intialize service: failed to create api client: %w", err)
	}
	return nil
}

func (s *service) Start() {
	for range time.Tick(s.interval) {
		ip, err := ipify.GetIP()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to get IP")
			continue
		}
		s.logger.Info().Str("IP", ip).Msg("ip to update")
		record := dnsapi.Record{
			Name:    s.config.RecordName,
			Content: ip,
			TTL:     s.config.TTL,
		}
		if strings.Contains(ip, ":") {
			record.Type = "AAAA"
		} else {
			record.Type = "A"
		}
		if (*s.client).Exists(s.config.ZoneID, record) {
			if err := (*s.client).Update(s.config.ZoneID, record); err != nil {
				s.logger.Error().Err(err).Msg("failed to update record")
				continue
			}
			s.logger.Info().Msg("updated record successfully")
		} else {
			s.logger.Info().Msgf("record %s seems not exists, try to make new one...", s.config.RecordName)
			if err := (*s.client).Create(s.config.ZoneID, record); err != nil {
				s.logger.Error().Err(err).Msg("failed to create record")
				continue
			}
			s.logger.Info().Msg("created record successfully")
		}
	}
}
