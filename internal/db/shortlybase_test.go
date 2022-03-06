//go:build unit
// +build unit

package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/burkaydurdu/shortly/pkg/log"

	"github.com/burkaydurdu/shortly/config"
	"github.com/stretchr/testify/suite"
)

type ShortlyBaseSuite struct {
	suite.Suite
	db          *DB
	shortlyBase *ShortlyBase
}

func (s *ShortlyBaseSuite) SetupSuite() {
	var db = new(DB)
	s.db = db

	conf, _ := config.New()
	conf.MemoryPath = "../../.mem"
	conf.MemoryFileName = "test_shortly"

	var shortlyLog = log.ShortlyLog{
		Tag: "TEST",
	}

	s.shortlyBase = &ShortlyBase{
		Log:    &shortlyLog,
		Config: conf,
	}
}

func (s *ShortlyBaseSuite) SetupTest() {
	removeDB(s.shortlyBase.Config)
}

func (s *ShortlyBaseSuite) TearDownSuite() {
	removeDB(s.shortlyBase.Config)
}

func (s *ShortlyBaseSuite) TestFindByCode() {
	db, err := s.shortlyBase.InitialDB()

	db.ShortURL = append(db.ShortURL, Shortly{
		Code: "xxxxxx",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), db.FindByCode("xxxxxx"))
	assert.Nil(s.T(), db.FindByCode("123456"))
}

func (s *ShortlyBaseSuite) TestInitialDB() {
	db, err := s.shortlyBase.InitialDB()

	assert.Zero(s.T(), len(db.ShortURL))
	assert.NoError(s.T(), err)
}

func (s *ShortlyBaseSuite) TestDB() {
	code := "XXXXXX"

	s.db.ShortURL = append(s.db.ShortURL, Shortly{Code: code})

	_ = s.shortlyBase.SaveToFile(s.db)

	newDB, _ := s.shortlyBase.ReadFromFile()

	assert.Equal(s.T(), len(newDB.ShortURL), 1)
	assert.Equal(s.T(), newDB.ShortURL[0].Code, code)
}

func removeDB(config *config.Config) {
	_ = os.Remove(config.MemoryPath + "/" + config.MemoryFileName + ".json")
}

// Run Test Suit
func TestShortlyBaseSuite(t *testing.T) {
	suite.Run(t, new(ShortlyBaseSuite))
}
