package options

import (
	"github.com/KusionStack/karbour/pkg/apiserver"
	"github.com/spf13/pflag"
)

type ElasticSearchConfig struct {
	Addresses []string
	UserName  string
	Password  string
}

type SearchStorageOptions struct {
	SearchStorageType      string
	ElasticSearchAddresses []string
	ElasticSearchName      string
	ElasticSearchPassword  string
}

func NewSearchStorageOptions() *SearchStorageOptions {
	return &SearchStorageOptions{}
}

func (o *SearchStorageOptions) Validate() []error {
	return nil
}

func (o *SearchStorageOptions) ApplyTo(config *apiserver.ExtraConfig) error {
	config.SearchStorageType = o.SearchStorageType
	config.ElasticSearchAddresses = o.ElasticSearchAddresses
	config.ElasticSearchName = o.ElasticSearchName
	config.ElasticSearchPassword = o.ElasticSearchPassword
	return nil
}

// AddFlags adds flags for a specific Option to the specified FlagSet
func (o *SearchStorageOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.SearchStorageType, "search-storage-type", "", "The search storage type")
	fs.StringSliceVar(&o.ElasticSearchAddresses, "elastic-search-addresses", nil, "The elastic search address")
	fs.StringVar(&o.ElasticSearchName, "elastic-search-username", "", "The elastic search username")
	fs.StringVar(&o.ElasticSearchPassword, "elastic-search-password", "", "The elastic search password")
}
