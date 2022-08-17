package network

type Config struct {
    IPv4 struct {
        Address    string
        Gateway    string
        RouteTable int
        TProxyPort int
    }
    V4RouteTable int
    V6RouteTable int
    V4TProxyPort int
    V6TProxyPort int
}

func Load() {

}
