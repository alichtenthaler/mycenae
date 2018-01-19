package collector

import (
	"fmt"
	"time"

	"github.com/uol/gobol"
	"github.com/uol/mycenae/lib/utils"
	"strconv"
)

func (collector *Collector) makePacket(packet *Point, rcvMsg TSDBpoint, number bool) gobol.Error {

	if number {
		if rcvMsg.Value == nil {
			return errValidation(`Wrong Format: Field "value" is required. NO information will be saved`)
		}
	} else {
		if rcvMsg.Text == "" {
			return errValidation(`Wrong Format: Field "text" is required. NO information will be saved`)
		}

		if len(rcvMsg.Text) > 10000 {
			return errValidation(`Wrong Format: Field "text" can not have more than 10k`)
		}
	}

	lt := len(rcvMsg.Tags)

	if lt == 0 {
		return errValidation(`Wrong Format: At least one tag is required. NO information will be saved`)
	}

	if !collector.validKey.MatchString(rcvMsg.Metric) {
		return errValidation(
			fmt.Sprintf(
				`Wrong Format: Field "metric" (%s) is not well formed. NO information will be saved`,
				rcvMsg.Metric,
			),
		)
	}

	if keyset, ok := rcvMsg.Tags["ksid"]; !ok {
		return errValidation(`Wrong Format: Tag "ksid" is required. NO information will be saved`)
	} else {
		packet.Keyset = keyset
	}

	if !collector.keySet.IsKeySetNameValid(packet.Keyset) {
		return errValidation(
			fmt.Sprintf(
				`Wrong Format: Field "ksid" (%s) is not well formed. NO information will be saved`,
				packet.Keyset,
			),
		)
	}

	keySetExists, gerr := collector.keySet.KeySetExists(packet.Keyset)
	if gerr != nil {
		return errISE("makePacket", "error checking keyspace existence", gerr)
	}
	if !keySetExists {
		return errValidation(`"ksid" not exists. NO information will be saved`)
	}

	if strTTL, ok := rcvMsg.Tags["ttl"]; !ok {
		packet.TTL = collector.settings.DefaultTTL
		lt++
	} else {
		ttl64, err := strconv.ParseUint(strTTL, 10, 8)
		if err != nil {
			return errValidation(`Wrong Format: Tag "ttl" must be a positive number. NO information will be saved`)
		}
		ttl := uint8(ttl64)
		if _, ok := collector.keyspaceTTLMap[ttl]; !ok {
			ttl = collector.settings.DefaultTTL
		}
		packet.TTL = ttl
	}
	rcvMsg.Tags["ttl"] = strconv.Itoa(int(packet.TTL))

	if lt == 2 {
		return errValidation(`Wrong Format: At least one tag other than "ksid" and "ttl" is required. NO information will be saved`)
	}

	for k, v := range rcvMsg.Tags {
		if !collector.validKey.MatchString(k) {
			return errValidation(
				fmt.Sprintf(
					`Wrong Format: Tag key (%s) is not well formed. NO information will be saved`,
					k,
				),
			)
		}
		if !collector.validKey.MatchString(v) {
			return errValidation(
				fmt.Sprintf(
					`Wrong Format: Tag value (%s) is not well formed. NO information will be saved`,
					v,
				),
			)
		}
	}

	if rcvMsg.Timestamp == 0 {
		packet.Timestamp = utils.GetTimeNoMillis()
	} else {
		truncated, err := utils.MilliToSeconds(rcvMsg.Timestamp)
		if err != nil {
			return errBR("makePacket", err.Error(), err)
		}
		packet.Timestamp = truncated
	}

	packet.Number = number
	packet.Message = rcvMsg
	packet.ID = GenerateID(rcvMsg)
	if !number {
		packet.ID = fmt.Sprintf("T%v", packet.ID)
	}
	year, week := time.Unix(0, packet.Timestamp*1e+6).ISOWeek()
	packet.Bucket = fmt.Sprintf("%v%v", year, week)

	return nil
}