package common

//------------------------------------------------------------------------------

// NetworkConfiguration holds the network configuration information
type NetworkConfiguration struct {
  Name        string `yaml:"Name"`        // URL of API endpoint
  IPv4CIDR    string `yaml:"IPv4CIDR"`    // IPv4 CIDR
  IPv4Gateway string `yaml:"IPv4Gateway"` // IPv4 gateway address
  IPv4Start   string `yaml:"IPv4Start"`   // IPv4 DHCP pool start
  IPv4Stop    string `yaml:"IPv4Stop"`    // IPv4 DHCP pool end
}

//------------------------------------------------------------------------------

// NetworkEndpoint holds the network endpoint information
type NetworkEndpoint struct {
  Name        string `yaml:"Name"`        // URL of API endpoint
  UUID1       string `yaml:"UUID"`        // unique universal identifier of network
  UUID2       string `yaml:"UUID"`        // unique universal identifier of subnet
  IPv4CIDR    string `yaml:"IPv4CIDR"`    // IPv4 CIDR
  IPv4Gateway string `yaml:"IPv4Gateway"` // IPv4 gateway address
  IPv4Start   string `yaml:"IPv4Start"`   // IPv4 DHCP pool start
  IPv4Stop    string `yaml:"IPv4Stop"`    // IPv4 DHCP pool end
}

//------------------------------------------------------------------------------
