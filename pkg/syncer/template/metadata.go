package template

// Config is the struct that defines the metadata of the agent template.
type Config struct {
	ExternalEndpoint string
	CaCert           string
	CaKey            string
	StorageAddresses []string
	ClusterName      string
	ClusterMode      string

	Level int
	Mode  string
}
