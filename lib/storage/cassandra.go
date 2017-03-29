package storage

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gocql/gocql"
	"github.com/uol/gobol"
)

type Cassandra struct {
	session            *gocql.Session
	writeConsistencies []gocql.Consistency
	readConsistencies  []gocql.Consistency
}

type cassPoints struct {
	date int64
	blob []byte
}

func (cass *Cassandra) SetWriteConsistencies(consistencies []gocql.Consistency) {
	cass.writeConsistencies = consistencies
}

func (cass *Cassandra) SetReadConsistencies(consistencies []gocql.Consistency) {
	cass.readConsistencies = consistencies
}

func (cass *Cassandra) InsertBucket(ksid, tsid string, timestamp int64, value []byte) gobol.Error {
	start := time.Now()

	year, week := time.Unix(timestamp, 0).ISOWeek()
	tsid = fmt.Sprintf("%v%v%v", year, week, tsid)

	var err error
	for _, cons := range cass.writeConsistencies {
		if err = cass.session.Query(
			fmt.Sprintf(`INSERT INTO %v.timeserie (id, date , value) VALUES (?, ?, ?)`, ksid),
			tsid,
			timestamp,
			value,
		).Consistency(cons).RoutingKey([]byte(tsid)).Exec(); err != nil {
			statsInsertFBerror(ksid, "ts_number_stamp")
			gblog.WithFields(
				logrus.Fields{
					"package": "collector/persistence",
					"func":    "insertPoint",
				},
			).Error(err)
			continue
		}
		statsInsert(ksid, "ts_number_stamp", time.Since(start))
		return nil
	}
	statsInsertQerror(ksid, "ts_number_stamp")
	return errPersist("InsertPoint", err)
}

func (cass *Cassandra) ReadBucket(ksid, tsid string, start, end int64) ([]cassPoints, int, gobol.Error) {
	track := time.Now()
	start = start - 3600
	end++

	var date int64
	var value []byte
	var err error

	buckets := []string{}

	w := start

	for {

		year, week := time.Unix(w, 0).ISOWeek()

		buckets = append(buckets, fmt.Sprintf("%v%v", year, week))

		if w > end {
			break
		}

		w += 604800

	}

	//year, week := time.Unix(start, 0).ISOWeek()
	//tsid = fmt.Sprintf("%v%v%v", year, week, tsid)

	for _, cons := range cass.readConsistencies {
		iter := cass.session.Query(
			fmt.Sprintf(
				`SELECT date, value FROM %v.timeserie WHERE id= ? AND date > ? AND date < ? ALLOW FILTERING`,
				ksid,
			),
			tsid,
			start,
			end,
		).Consistency(cons).RoutingKey([]byte(tsid)).Iter()

		points := []cassPoints{}
		var count int

		for iter.Scan(&date, &value) {
			count++
			point := cassPoints{
				date: date,
				blob: value,
			}
			points = append(points, point)
		}

		if err = iter.Close(); err != nil {

			gblog.WithFields(logrus.Fields{
				"package": "storage/cassandra",
				"func":    "getTStamp",
			}).Error(err)

			if err == gocql.ErrNotFound {
				statsSelect(ksid, "ts_number_stamp", time.Since(track))
				return []cassPoints{}, 0, errNoContent("getTStamp")
			}

			statsSelectQerror(ksid, "ts_number_stamp")
			continue
		}
		statsSelect(ksid, "ts_number_stamp", time.Since(track))
		return points, count, nil
	}
	statsSelectFerror(ksid, "ts_number_stamp")
	return []cassPoints{}, 0, errPersist("getTStamp", err)
}

func (cass *Cassandra) InsertText(ksid, tsid string, timestamp int64, text string) gobol.Error {
	start := time.Now()
	var err error
	for _, cons := range cass.writeConsistencies {
		if err = cass.session.Query(
			fmt.Sprintf(`INSERT INTO %v.ts_text_stamp (id, date , value) VALUES (?, ?, ?)`, ksid),
			tsid,
			timestamp,
			text,
		).Consistency(cons).RoutingKey([]byte(tsid)).Exec(); err != nil {
			statsInsertQerror(ksid, "ts_text_stamp")
			gblog.WithFields(
				logrus.Fields{
					"package": "collector/persistence",
					"func":    "InsertText",
				},
			).Error(err)
			continue
		}
		statsInsert(ksid, "ts_text_stamp", time.Since(start))
		return nil
	}
	statsInsertFBerror(ksid, "ts_text_stamp")
	return errPersist("InsertText", err)
}

func (cass *Cassandra) InsertError(id, msg, errMsg string, date time.Time) gobol.Error {
	start := time.Now()
	var err error
	for _, cons := range cass.writeConsistencies {
		if err = cass.session.Query(
			`INSERT INTO ts_error (code, tsid, error, message, date) VALUES (?, ?, ?, ?, ?)`,
			0,
			id,
			errMsg,
			msg,
			date,
		).Consistency(cons).RoutingKey([]byte(id)).Exec(); err != nil {
			statsInsertQerror("default", "ts_error")
			gblog.WithFields(
				logrus.Fields{
					"package": "collector/persistence",
					"func":    "InsertError",
				},
			).Error(err)
			continue
		}
		statsInsert("default", "ts_error", time.Since(start))
		return nil
	}
	statsInsertFBerror("default", "ts_error")
	return errPersist("InsertError", err)
}
