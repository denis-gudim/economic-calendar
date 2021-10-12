package loading

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	"economic-calendar/loader/investing"

	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type DictionariesLoaderService struct {
	investingRepository InvestingDataReciver
	countriesRepository CountriesDataReciver
	logger              *log.Logger
	config              app.Config
}

func NewDictionariesLoaderService(cnf app.Config, logger *log.Logger) *DictionariesLoaderService {

	return &DictionariesLoaderService{
		investingRepository: investing.NewInvestingRepository(cnf, logger),
		countriesRepository: data.NewCountriesRepository(cnf),
		logger:              logger,
		config:              cnf,
	}
}

func (s *DictionariesLoaderService) Load() error {

	fmtError := func(err error) error {
		return xerrors.Errorf("countries dictionary loading failed: %w", err)
	}

	s.logger.Info("countries dictionary loading started...")

	countries, err := s.countriesRepository.GetAll()

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

	invCountries, err := s.investingRepository.GetCountries()

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

		err = s.countriesRepository.Save(c)

		if err != nil {
			return fmtError(err)
		}

	}

	s.logger.Info("countries dictionary loading finished")

	return nil
}
