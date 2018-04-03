package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/uol/gobol/cassandra"
	"github.com/uol/gobol/loader"
	"github.com/uol/gobol/rubber"
	"github.com/uol/gobol/saw"
	"github.com/uol/gobol/snitch"
	"go.uber.org/zap"

	"github.com/uol/mycenae/lib/cache"
	"github.com/uol/mycenae/lib/collector"
	"github.com/uol/mycenae/lib/keyset"
	"github.com/uol/mycenae/lib/keyspace"
	"github.com/uol/mycenae/lib/memcached"
	"github.com/uol/mycenae/lib/metadata"
	"github.com/uol/mycenae/lib/persistence"
	"github.com/uol/mycenae/lib/plot"
	"github.com/uol/mycenae/lib/rest"
	"github.com/uol/mycenae/lib/structs"
	"github.com/uol/mycenae/lib/tsstats"
	"github.com/uol/mycenae/lib/udp"
)

func main() {
	fmt.Println("Starting Mycenae")
	//Parse of command line arguments.
	var confPath string
	var devMode bool

	flag.StringVar(&confPath, "config", "config.toml", "path to configuration file")
	flag.BoolVar(&devMode, "devMode", false, "enable/disable dev mode (all ttls are configured to one day)")
	flag.Parse()

	//Load conf file.
	settings := new(structs.Settings)

	err := loader.ConfToml(confPath, &settings)
	if err != nil {
		log.Fatalln("ERROR - Loading Config file: ", err)
	} else {
		fmt.Println("config file loaded.")
	}

	tsLogger := new(structs.TsLog)
	tsLogger.General, err = saw.New(settings.Logs.General.Level, settings.Logs.Environment)
	if err != nil {
		log.Fatalln("ERROR - Starting logger: ", err)
	}
	tsLogger.General = tsLogger.General.With(zap.String("type", settings.Logs.General.Prefix))

	tsLogger.Stats, err = saw.New(settings.Logs.Stats.Level, settings.Logs.Environment)
	if err != nil {
		log.Fatalln("ERROR - Starting logger: ", err)
	}
	tsLogger.Stats = tsLogger.Stats.With(zap.String("type", settings.Logs.Stats.Prefix))

	lf := []zapcore.Field{
		zap.String("package", "main"),
		zap.String("func", "main"),
	}

	if devMode {
		tsLogger.General.Info("DEV MODE IS ENABLED!", lf...)
	}

	go func() {
		err := http.ListenAndServe("0.0.0.0:6666", nil)

		if err != nil {
			tsLogger.General.Error(err.Error(), lf...)
		}
	}()

	sts, err := snitch.New(tsLogger.Stats, settings.Stats)
	if err != nil {
		tsLogger.General.Fatal(fmt.Sprintf("ERROR - Starting stats: %s", err.Error()), lf...)
	}

	tssts, err := tsstats.New(tsLogger.General, sts, settings.Stats.Interval, settings.Stats.KSID, settings.Stats.Tags["ttl"])
	if err != nil {
		tsLogger.General.Error(err.Error(), lf...)
		os.Exit(1)
	}

	cass, err := cassandra.New(settings.Cassandra)
	if err != nil {
		tsLogger.General.Fatal(fmt.Sprintf("ERROR - Connecting to cassandra: %s", err.Error()), lf...)
		os.Exit(1)
	}
	defer cass.Close()

	// --- Including metadata and persistence ---
	meta, err := metadata.Create(
		settings.ElasticSearch.Cluster,
		tsLogger.General,
		tssts,
	)
	if err != nil {
		tsLogger.General.Fatal(fmt.Sprintf("error creating metadata backend: %s", err.Error()), lf...)
		os.Exit(1)
	}

	storage, err := persistence.NewStorage(
		settings.Cassandra.Keyspace,
		settings.Cassandra.Username,
		tsLogger.General,
		cass,
		meta,
		tssts,
		devMode,
		settings.DefaultTTL,
	)
	if err != nil {
		tsLogger.General.Fatal(fmt.Sprintf("error creating persistence backend: %s", err.Error()), lf...)
		os.Exit(1)
	}
	// --- End of metadata and persistence ---

	es, err := rubber.New(tsLogger.General, settings.ElasticSearch.Cluster)
	if err != nil {
		tsLogger.General.Fatal(fmt.Sprintf("ERROR - Connecting to elasticsearch: %s", err.Error()), lf...)
		os.Exit(1)
	}

	ks := keyspace.New(
		tssts,
		storage,
		devMode,
		settings.DefaultTTL,
		settings.MaxAllowedTTL,
	)

	jsonStr, _ := json.Marshal(settings.DefaultKeyspaces)
	tsLogger.General.Info(fmt.Sprintf("creating default keyspaces: %s", jsonStr), lf...)
	keyspaceTTLMap := map[uint8]string{}
	for k, ttl := range settings.DefaultKeyspaces {
		gerr := ks.Storage.CreateKeyspace(k,
			settings.DefaultKeyspaceData.Datacenter,
			settings.DefaultKeyspaceData.Contact,
			settings.DefaultKeyspaceData.ReplicationFactor,
			ttl)
		keyspaceTTLMap[ttl] = k
		if gerr != nil && gerr.StatusCode() != http.StatusConflict {
			tsLogger.General.Fatal(fmt.Sprintf("error creating kayspace '%s': %s", k, gerr.Message()), lf...)
			os.Exit(1)
		}
	}

	mc, err := memcached.New(tssts, &settings.Memcached)
	if err != nil {
		tsLogger.General.Fatal(err.Error(), lf...)
		os.Exit(1)
	}

	kc := cache.NewKeyspaceCache(mc, ks)
	keySet := keyset.NewKeySet(es, tssts, mc)

	jsonStr, _ = json.Marshal(settings.DefaultKeysets)
	tsLogger.General.Info(fmt.Sprintf("creating default keysets: %s", jsonStr), lf...)
	for _, v := range settings.DefaultKeysets {
		exists, err := keySet.KeySetExists(v)
		if err != nil {
			tsLogger.General.Fatal(fmt.Sprintf("error checking keyset '%s' existence: %s", v, err.Error()), lf...)
			os.Exit(1)
		}
		if !exists {
			tsLogger.General.Info(fmt.Sprintf("creating default keyset '%s'", v), lf...)
			err = keySet.CreateIndex(v)
			if err != nil {
				tsLogger.General.Fatal(fmt.Sprintf("error creating keyset '%s': %s", v, err.Error()), lf...)
				os.Exit(1)
			}
		}
	}

	coll, err := collector.New(tsLogger, tssts, cass, es, kc, settings, keyspaceTTLMap, keySet)
	if err != nil {
		log.Println(err)
		return
	}

	uV2server := udp.New(tsLogger.General, settings.UDPserverV2, coll)
	uV2server.Start()

	collectorV1 := collector.UDPv1{}
	uV1server := udp.New(tsLogger.General, settings.UDPserver, collectorV1)
	uV1server.Start()

	p, err := plot.New(
		tsLogger.General,
		tssts,
		cass,
		es,
		kc,
		settings.ElasticSearch.Index,
		settings.MaxTimeseries,
		settings.MaxConcurrentTimeseries,
		settings.MaxConcurrentReads,
		settings.LogQueryTSthreshold,
		keyspaceTTLMap,
		keySet,
		settings.DefaultTTL,
	)
	if err != nil {
		tsLogger.General.Fatal(err.Error(), lf...)
		os.Exit(1)
	}

	tsRest := rest.New(
		tsLogger,
		sts,
		p,
		ks,
		mc,
		coll,
		settings.HTTPserver,
		settings.Probe.Threshold,
		keySet,
	)
	tsRest.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	fmt.Println("mycenae started successfully")
	for {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			stop(tsLogger, tsRest, coll)
			return
		case syscall.SIGHUP:
			//THIS IS A HACK DO NOT EXTEND IT. THE FEATURE IS NICE BUT NEEDS TO BE DONE CORRECTLY!!!!!
			settings := new(structs.Settings)
			var err error

			if strings.HasSuffix(confPath, ".json") {
				err = loader.ConfJson(confPath, &settings)
			} else if strings.HasSuffix(confPath, ".toml") {
				err = loader.ConfToml(confPath, &settings)
			}
			if err != nil {
				tsLogger.General.Error(fmt.Sprintf("ERROR - Loading Config file: %s", err.Error()), lf...)
				continue
			} else {
				tsLogger.General.Info("config file loaded", lf...)
			}
		}
	}
}

func stop(logger *structs.TsLog, rest *rest.REST, collector *collector.Collector) {

	lf := []zapcore.Field{
		zap.String("package", "main"),
		zap.String("func", "stop"),
	}

	logger.General.Info("stopping REST", lf...)
	rest.Stop()
	logger.General.Info("REST stopped", lf...)

	logger.General.Info("stopping UDPv2", lf...)
	collector.Stop()
	logger.General.Info("UDPv2 stopped", lf...)
}
