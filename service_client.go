package bbs

import (
	"path"
	"time"

	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/cloudfoundry-incubator/consuladapter"
	"github.com/cloudfoundry-incubator/locket"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
)

const (
	CellSchemaKey    = "cell"
	BBSLockSchemaKey = "bbs_lock"
)

func CellSchemaRoot() string {
	return locket.LockSchemaPath(CellSchemaKey)
}

func CellSchemaPath(cellID string) string {
	return locket.LockSchemaPath(CellSchemaKey, cellID)
}

func BBSLockSchemaPath() string {
	return locket.LockSchemaPath(BBSLockSchemaKey)
}

//go:generate counterfeiter -o fake_bbs/fake_service_client.go . ServiceClient

type ServiceClient interface {
	CellById(logger lager.Logger, cellId string) (*models.CellPresence, error)
	Cells(logger lager.Logger) (models.CellSet, error)
	CellEvents(logger lager.Logger) <-chan models.CellEvent
	NewCellPresenceRunner(logger lager.Logger, cellPresence *models.CellPresence, retryInterval, lockTTL time.Duration) ifrit.Runner
	NewBBSLockRunner(logger lager.Logger, bbsPresence *models.BBSPresence, retryInterval time.Duration) (ifrit.Runner, error)
	CurrentBBS(logger lager.Logger) (*models.BBSPresence, error)
	CurrentBBSURL(logger lager.Logger) (string, error)
}

type serviceClient struct {
	session      *consuladapter.Session
	consulClient consuladapter.Client
	clock        clock.Clock
}

func NewServiceClient(logger lager.Logger, client consuladapter.Client, lockTTL time.Duration, clock clock.Clock) ServiceClient {
	uuid, err := uuid.NewV4()
	if err != nil {
		logger.Fatal("construct-uuid-failed", err)
	}

	session, err := consuladapter.NewSessionNoChecks(uuid.String(), lockTTL, client)
	if err != nil {
		logger.Fatal("consul-session-failed", err)
	}

	return &serviceClient{
		session:      session,
		consulClient: client,
		clock:        clock,
	}
}

func (db *serviceClient) NewCellPresenceRunner(logger lager.Logger, cellPresence *models.CellPresence, retryInterval time.Duration, lockTTL time.Duration) ifrit.Runner {
	payload, err := models.ToJSON(cellPresence)
	if err != nil {
		panic(err)
	}

	return locket.NewPresence(logger, db.consulClient, CellSchemaPath(cellPresence.CellId), payload, db.clock, retryInterval, lockTTL)
}

func (db *serviceClient) Cells(logger lager.Logger) (models.CellSet, error) {
	kvPairs, _, err := db.consulClient.KV().List(CellSchemaRoot(), nil)
	if err != nil {
		bbsErr := models.ConvertError(convertConsulError(err))
		if bbsErr.Type != models.Error_ResourceNotFound {
			return nil, bbsErr
		}
	}

	if kvPairs == nil {
		err = consuladapter.NewPrefixNotFoundError(CellSchemaRoot())
		bbsErr := models.ConvertError(convertConsulError(err))
		if bbsErr.Type != models.Error_ResourceNotFound {
			return nil, bbsErr
		}
	}

	cellPresences := models.NewCellSet()
	for _, kvPair := range kvPairs {
		if kvPair.Session == "" {
			continue
		}

		cell := kvPair.Value
		presence := new(models.CellPresence)
		err := models.FromJSON(cell, presence)
		if err != nil {
			logger.Error("failed-to-unmarshal-cells-json", err)
			continue
		}
		cellPresences.Add(presence)
	}

	return cellPresences, nil
}

func (db *serviceClient) CellById(logger lager.Logger, cellId string) (*models.CellPresence, error) {
	value, err := db.session.GetAcquiredValue(CellSchemaPath(cellId))
	if err != nil {
		return nil, convertConsulError(err)
	}

	presence := new(models.CellPresence)
	err = models.FromJSON(value, presence)
	if err != nil {
		return nil, models.NewError(models.Error_InvalidJSON, err.Error())
	}

	return presence, nil
}

func (db *serviceClient) CellEvents(logger lager.Logger) <-chan models.CellEvent {
	logger = logger.Session("cell-events")

	events := make(chan models.CellEvent)
	go func() {
		disappeared := db.session.WatchForDisappearancesUnder(logger, CellSchemaRoot())

		for {
			select {
			case keys, ok := <-disappeared:
				if !ok {
					return
				}

				cellIDs := make([]string, len(keys))
				for i, key := range keys {
					cellIDs[i] = path.Base(key)
				}
				logger.Info("cell-disappeared", lager.Data{"cell-ids": cellIDs})
				events <- models.NewCellDisappearedEvent(cellIDs)
			}
		}
	}()

	return events
}

func (db *serviceClient) NewBBSLockRunner(logger lager.Logger, bbsPresence *models.BBSPresence, retryInterval time.Duration) (ifrit.Runner, error) {
	bbsPresenceJSON, err := models.ToJSON(bbsPresence)
	if err != nil {
		return nil, err
	}
	return locket.NewLock(db.session, locket.LockSchemaPath("bbs_lock"), bbsPresenceJSON, db.clock, retryInterval, logger), nil
}

func (db *serviceClient) CurrentBBS(logger lager.Logger) (*models.BBSPresence, error) {
	value, err := db.session.GetAcquiredValue(BBSLockSchemaPath())
	if err != nil {
		return nil, convertConsulError(err)
	}

	presence := new(models.BBSPresence)
	err = models.FromJSON(value, presence)
	if err != nil {
		return nil, err
	}

	return presence, nil
}

func (db *serviceClient) CurrentBBSURL(logger lager.Logger) (string, error) {
	presence, err := db.CurrentBBS(logger)
	if err != nil {
		return "", err
	}

	return presence.URL, nil
}

func convertConsulError(err error) error {
	switch err.(type) {
	case consuladapter.KeyNotFoundError:
		return models.NewError(models.Error_ResourceNotFound, err.Error())
	case consuladapter.PrefixNotFoundError:
		return models.NewError(models.Error_ResourceNotFound, err.Error())
	default:
		return models.NewError(models.Error_UnknownError, err.Error())
	}
}