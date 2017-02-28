package collector

import (
	"time"
)

func statsUDPv1() {
	go statsIncrement(
		"points.received",
		map[string]string{"protocol": "udp", "api": "v1"},
	)
}

func statsUDP(ks, vt string) {
	go statsIncrement(
		"points.received",
		map[string]string{"protocol": "udp", "api": "v2", "keyspace": ks, "type": vt},
	)
}

func statsProcTime(ks string, d time.Duration) {
	go statsValueAdd(
		"points.processes_time",
		map[string]string{"keyspace": ks},
		float64(d.Nanoseconds())/float64(time.Millisecond),
	)
}

func statsLostMeta() {
	go statsIncrement(
		"meta.lost",
		map[string]string{},
	)
}

func statsUDPerror(ks, vt string) {
	go statsIncrement(
		"points.received.error",
		map[string]string{"protocol": "udp", "api": "v2", "keyspace": ks, "type": vt},
	)
}

func statsInsertQerror(ks, cf string) {
	go statsIncrement(
		"cassandra.query.error",
		map[string]string{"keyspace": ks, "column_family": cf, "operation": "insert"},
	)
}

func statsInsertFBerror(ks, cf string) {
	go statsIncrement(
		"cassandra.fallback.error",
		map[string]string{"keyspace": ks, "column_family": cf, "operation": "insert"},
	)
}

func statsIndexError(i, t, m string) {
	tags := map[string]string{"method": m}
	if i != "" {
		tags["index"] = i
	}
	if t != "" {
		tags["type"] = t
	}
	go statsIncrement("elastic.request.error", tags)
}

func statsIndex(i, t, m string, d time.Duration) {
	tags := map[string]string{"method": m}
	if i != "" {
		tags["index"] = i
	}
	if t != "" {
		tags["type"] = t
	}
	go statsIncrement("elastic.request", tags)
	go statsValueAdd(
		"elastic.request.duration",
		tags,
		float64(d.Nanoseconds())/float64(time.Millisecond),
	)
}

func statsBulkPoints() {
	go statsIncrement("elastic.bulk.points", map[string]string{})
}

func statsInsert(ks, cf string, d time.Duration) {
	go statsIncrement("cassandra.query", map[string]string{"keyspace": ks, "column_family": cf, "operation": "insert"})
	go statsValueAdd(
		"cassandra.query.duration",
		map[string]string{"keyspace": ks, "column_family": cf, "operation": "insert"},
		float64(d.Nanoseconds())/float64(time.Millisecond),
	)
}

func statsPoints(ks, vt string) {
	go statsIncrement(
		"points.received",
		map[string]string{"protocol": "http", "api": "v2", "keyspace": ks, "type": vt},
	)
}

func statsPointsError(ks, vt string) {
	go statsIncrement(
		"points.received.error",
		map[string]string{"protocol": "http", "api": "v2", "keyspace": ks, "type": vt},
	)
}

func statsIncrement(metric string, tags map[string]string) {
	stats.Increment("collector", metric, tags)
}

func statsValueAdd(metric string, tags map[string]string, v float64) {
	stats.ValueAdd("collector", metric, tags, v)
}
