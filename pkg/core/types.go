package core

import (
	"fmt"
	"strings"
)

type Locator struct {
	Cluster    string `json:"cluster" yaml:"cluster"`
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Group      string `json:"group" yaml:"group"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	Name       string `json:"name" yaml:"name"`
}

func (c *Locator) ToSQL() string {
	var conditions []string

	if c.Cluster != "" {
		conditions = append(conditions, fmt.Sprintf("cluster='%s'", c.Cluster))
	}
	if c.ApiVersion != "" {
		conditions = append(conditions, fmt.Sprintf("apiVersion='%s'", c.ApiVersion))
	}
	if c.Group != "" {
		conditions = append(conditions, fmt.Sprintf("group='%s'", c.Group))
	}
	if c.Namespace != "" {
		conditions = append(conditions, fmt.Sprintf("namespace='%s'", c.Namespace))
	}
	if c.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name='%s'", c.Name))
	}

	if len(conditions) > 0 {
		return "SELECT * from resources WHERE " + strings.Join(conditions, " AND ")
	} else {
		return "SELECT * from resources"
	}
}
