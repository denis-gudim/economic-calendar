package loading

import (
	"context"

	"github.com/denis-gudim/economic-calendar/loader"
	"github.com/denis-gudim/economic-calendar/loader/data"

	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type DictionariesLoaderService struct {
	investingRepository InvestingDataReciver
	countriesRepository CountriesDataReciver
	logger              *log.Logger
	config              *loader.Config
}

func NewDictionariesLoaderService(cnf *loader.Config,
	logger *log.Logger,
	investingRepository InvestingDataReciver,
	countriesRepository CountriesDataReciver) *DictionariesLoaderService {

	return &DictionariesLoaderService{
		investingRepository: investingRepository,
		countriesRepository: countriesRepository,
		logger:              logger,
		config:              cnf,
	}
}

func (s *DictionariesLoaderService) Load(ctx context.Context) error {

	fmtError := func(err error) error {
		return xerrors.Errorf("countries dictionary loading failed: %w", err)
	}

	s.logger.Info("countries dictionary loading started...")

	countries, err := s.countriesRepository.GetAll(ctx)

	if err != nil {
		return fmtError(err)
	}

	load := false

	for _, c := range countries {
		if len(c.NameTranslations) == 0 {
			load = true
			break
		}
	}

	if !load {
		s.logger.Info("countries dictionary loading skiped")
		return nil
	}

	invCountries, err := s.investingRepository.GetCountries(ctx)

	if err != nil {
		return fmtError(err)
	}

	for _, c := range countries {
		ic, ok := invCountries[c.Id]

		if !ok {
			continue
		}

		c.NameTranslations = data.Translations{}

		for _, icl := range ic {
			c.NameTranslations[icl.LanguageId] = icl.Title
		}

		err = s.countriesRepository.Save(ctx, c)

		if err != nil {
			return fmtError(err)
		}

	}

	s.logger.Info("countries dictionary loading finished")

	return nil
}
