package elastic

import (
	"github.com/kubemq-hub/builder/connector/common"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.elasticsearch").
		SetDescription("Elastic Search Target").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("host").
				SetDescription("Set Elastic Search Urls").
				SetMust(true).
				SetDefault("http://localhost:9200"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("username").
				SetDescription("Set Elastic Search username").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("password").
				SetDescription("Set Elastic Search password").
				SetMust(false).
				SetDefault(""),
		).
		AddProperty(
			common.NewProperty().
				SetKind("bool").
				SetName("sniff").
				SetDescription("Set Elastic Search sniff mode").
				SetMust(false).
				SetDefault("true"),
		)
}
