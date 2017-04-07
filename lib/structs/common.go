package structs

import (
	"github.com/Sirupsen/logrus"
	"github.com/uol/gobol/cassandra"
	"github.com/uol/gobol/rubber"
	"github.com/uol/gobol/saw"
	"github.com/uol/gobol/snitch"
	"github.com/uol/mycenae/lib/cluster"
)

type TsLog struct {
	General *logrus.Logger
	Stats   *logrus.Logger
}

type SettingsHTTP struct {
	Path string
	Port string
	Bind string
}

type SettingsUDP struct {
	Port       string
	ReadBuffer int
}

type Settings struct {
	ReadConsistency         []string
	WriteConsisteny         []string
	BoltPath                string
	WALPath                 string
	MaxTimeseries           int
	MaxConcurrentTimeseries int
	MaxConcurrentReads      int
	LogQueryTSthreshold     int
	MaxConcurrentPoints     int
	MaxConcurrentBulks      int
	MaxMetaBulkSize         int
	MetaBufferSize          int
	MetaSaveInterval        string
	CompactionStrategy      string
	HTTPserver              SettingsHTTP
	UDPserver               SettingsUDP
	UDPserverV2             SettingsUDP
	Cassandra               cassandra.Settings
	Cluster                 cluster.Config
	TTL                     struct {
		Max int
	}
	Logs struct {
		General saw.Settings
		Stats   saw.Settings
	}
	Stats     snitch.Settings
	StatsFile struct {
		Path string
	}
	ElasticSearch struct {
		Cluster rubber.Settings
		Index   string
	}
	Probe struct {
		Threshold float64
	}
}
