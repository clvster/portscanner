package types

import "time"

type OpenPort struct {
	IP        string
	Port      int
	Proto     string
	Service   string
	Banner    string
	FirstSeen time.Time
	LastSeen  time.Time
}

func (p OpenPort) Key() string {
	return p.IP + ":" + itoa(p.Port) + "/" + p.Proto
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var b [20]byte

	i := len(b)

	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return string(b[i:])
}
