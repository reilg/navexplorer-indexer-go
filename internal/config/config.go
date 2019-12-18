package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Network            string
	Debug              bool
	Reindex            bool
	ReindexSize        uint
	BulkIndexSize      uint
	SoftForkBlockCycle uint
	SoftForkQuorum     uint
	Navcoind           NavcoindConfig
	ElasticSearch      ElasticSearchConfig
	ZeroMq             ZeroMqConfig
	DaoCfundConsensus  DaoCfundConsensusConfig
}

type NavcoindConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Ssl      bool
}

type ElasticSearchConfig struct {
	Hosts       []string
	Sniff       bool
	HealthCheck bool
	Debug       bool
	Username    string
	Password    string
	MappingDir  string
}

type ZeroMqConfig struct {
	Address string
}

type DaoCfundConsensusConfig struct {
	BlocksPerVotingCycle                uint
	Quorum                              uint
	MaxCountVotingCycleProposals        uint
	MaxCountVotingCyclePaymentRequests  uint
	VotesAcceptProposalPercentage       uint
	VotesRejectProposalPercentage       uint
	VotesAcceptPaymentRequestPercentage uint
	VotesRejectPaymentRequestPercentage uint
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Fatal("Unable to init config")
	}
	log.Info("Config init")
}

func Get() *Config {
	return &Config{
		Network:            getString("NAVCOIND_NETWORK", "mainnet"),
		SoftForkBlockCycle: getUint("SOFTFORK_BLOCKCYCLE", 20160),
		SoftForkQuorum:     getUint("SOFTFORK_QUORUM", 75),
		Debug:              getBool("DEBUG", false),
		Reindex:            getBool("REINDEX", false),
		ReindexSize:        getUint("REINDEX_SIZE", 200),
		BulkIndexSize:      getUint("BULK_INDEX_SIZE", 200),
		Navcoind: NavcoindConfig{
			Host:     getString("NAVCOIND_HOST", ""),
			Port:     getInt("NAVCOIND_PORT", 8332),
			User:     getString("NAVCOIND_USER", "user"),
			Password: getString("NAVCOIND_PASSWORD", "password"),
			Ssl:      getBool("NAVCOIND_SSL", false),
		},
		ElasticSearch: ElasticSearchConfig{
			Hosts:       getSlice("ELASTIC_SEARCH_HOSTS", make([]string, 0), ","),
			Sniff:       getBool("ELASTIC_SEARCH_SNIFF", true),
			HealthCheck: getBool("ELASTIC_SEARCH_HEALTH_CHECK", true),
			Debug:       getBool("ELASTIC_SEARCH_DEBUG", false),
			Username:    getString("ELASTIC_SEARCH_USERNAME", "/data/mappings"),
			Password:    getString("ELASTIC_SEARCH_PASSWORD", "/data/mappings"),
			MappingDir:  getString("ELASTIC_SEARCH_MAPPING_DIR", "/data/mappings"),
		},
		ZeroMq: ZeroMqConfig{
			Address: getString("ZEROMQ_ADDRESS", "tcp://navcoind:28332"),
		},
		DaoCfundConsensus: DaoCfundConsensusConfig{
			BlocksPerVotingCycle:                getUint("CFUND_BLOCKS_PER_VOTING_CYCLE", 20160),
			Quorum:                              getUint("CFUND_QUORUM", 50),
			MaxCountVotingCycleProposals:        getUint("CFUND_MAX_COUNT_VOTING_CYCLE_PROPOSALS", 6),
			MaxCountVotingCyclePaymentRequests:  getUint("CFUND_MAX_COUNT_VOTING_CYCLE_PAYMENT_REQUESTS", 8),
			VotesAcceptProposalPercentage:       getUint("CFUND_VOTES_ACCEPT_PROPOSAL_PERCENTAGE", 70),
			VotesRejectProposalPercentage:       getUint("CFUND_VOTES_REJECT_PROPOSAL_PERCENTAGE", 70),
			VotesAcceptPaymentRequestPercentage: getUint("CFUND_VOTES_ACCEPT_PAYMENT_REQUEST_PERCENTAGE", 70),
			VotesRejectPaymentRequestPercentage: getUint("CFUND_VOTES_REJECT_PAYMENT_REQUEST_PERCENTAGE", 70),
		},
	}
}

func getString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getInt(key string, defaultValue int) int {
	valStr := getString(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getUint(key string, defaultValue uint) uint {
	return uint(getInt(key, int(defaultValue)))
}

func getBool(key string, defaultValue bool) bool {
	valStr := getString(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

func getSlice(key string, defaultVal []string, sep string) []string {
	valStr := getString(key, "")
	if valStr == "" {
		return defaultVal
	}

	return strings.Split(valStr, sep)
}
