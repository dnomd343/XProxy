package proxy

type Config struct {
	Sniff         bool
	Redirect      bool
	V4TProxyPort  int
	V6TProxyPort  int
	LogLevel      string
	HttpInbounds  map[string]int
	SocksInbounds map[string]int
	AddOnInbounds []interface{}
}

type logObject struct {
	Loglevel string `json:"loglevel"`
	Access   string `json:"access"`
	Error    string `json:"error"`
}

type inboundsObject struct {
	Inbounds []interface{} `json:"inbounds"`
}

type sniffObject struct {
	Enabled      bool     `json:"enabled"`
	RouteOnly    bool     `json:"routeOnly"`
	DestOverride []string `json:"destOverride"`
}

type inboundObject struct {
	Tag            string      `json:"tag"`
	Port           int         `json:"port"`
	Protocol       string      `json:"protocol"`
	Settings       interface{} `json:"settings"`
	StreamSettings interface{} `json:"streamSettings"`
	Sniffing       sniffObject `json:"sniffing"`
}
