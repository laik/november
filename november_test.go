package november

import "testing"

import "strings"

type NewEnergyRaw struct {
	UUID    string `json:"UUID"`
	SESSION string `json:"SESSION"`
	VIN     string `json:"VIN"`
	CRC     int    `json:"CRC"`
	TIME    int    `json:"TIME"`
	CMD     int64  `json:"CMD"`
	ENC     int    `json:"ENC"`
	SEQ     int64  `json:"SEQ"`
	RSP     string `json:"RSP"`
	HEX     string `json:"RAW"`
}

func TestFormatRaw(t *testing.T) {
	data1 := "0d98fc50-eb03-418c-a40f-12883d7e138a	10.10.4.15:25431-T-117.136.43.131:4428	00000013789501687	-106	20171117195616	3	1	29718	-2	232303fe303030303030313337383935303136383701005d110b110b053201ffffffffff000b3f4cffffffffffff00ffffffff020101ffffffffffffffffffffff050006c2721001619e1e06ffffffffffffffffffffffffffff07000000000000000000080101ffffffffffff000100090101ffff96"
	//	data2 := "0d98fc50-eb03-418c-a40f-12883d7e138a	10.10.4.15:25431-T-117.136.43.131:4428	00000013789501687	-106	20171117195616	3	1	29718	-2"
	_split := func(s string) ([]string, error) {
		return strings.Split(s, "\t"), nil
	}
	t.Log(_split(data1))
	ner := new(NewEnergyRaw)
	XunmarshaText(ner, data1, _split)
	t.Log(ner)
}
