package enricher

import "encoding/xml"

type nmapRun struct {
	Hosts []nmapHost `xml:"host"`
}

type nmapHost struct {
	Addresses []nmapAddr `xml:"address"`
	Ports     []nmapPort `xml:"ports>port"`
}

type nmapAddr struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}

type nmapPort struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   int         `xml:"portid,attr"`
	Service  nmapService `xml:"service"`
}

type nmapService struct {
	Name    string   `xml:"name,attr"`
	Product string   `xml:"product,attr"`
	Version string   `xml:"version,attr"`
	CPEs    []string `xml:"cpe"`
}

type enriched struct {
	IP      string
	Port    int
	Proto   string
	Service string
	Product string
	Version string
	CPE     string
}

func parseNmapXML(data []byte) ([]enriched, error) {
	var run nmapRun
	if err := xml.Unmarshal(data, &run); err != nil {
		return nil, err
	}

	var result []enriched
	for _, h := range run.Hosts {
		ip := ""

		for _, a := range h.Addresses {
			if a.AddrType == "ipv4" || a.AddrType == "ipv6" {
				ip = a.Addr
				break
			}
		}

		if ip == "" {
			continue
		}

		for _, p := range h.Ports {
			proto := p.Protocol
			if proto == "" {
				proto = "tcp"
			}

			cpe := ""
			if len(p.Service.CPEs) > 0 {
				cpe = p.Service.CPEs[0]
			}

			result = append(result, enriched{
				IP:      ip,
				Port:    p.PortID,
				Proto:   proto,
				Service: p.Service.Name,
				Product: p.Service.Product,
				Version: p.Service.Version,
				CPE:     cpe,
			})
		}
	}

	return result, nil
}
