package parse

import "strconv"

func Port(port string) (p int) {
	if len(port) != 0 && port[0] == ':' {
		p, _ = strconv.Atoi(port[1:])
		return
	}
	p, _ = strconv.Atoi(port)
	return
}
