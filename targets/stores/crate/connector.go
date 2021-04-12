package crate

import (
	"github.com/kubemq-hub/builder/connector/common"
	"math"
)

func Connector() *common.Connector {
	return common.NewConnector().
		SetKind("stores.crate").
		SetDescription("Crate Target").
		SetName("Crate").
		SetProvider("").
		SetCategory("Store").
		SetTags("db","sql").
		AddProperty(
			common.NewProperty().
				SetKind("string").
				SetName("connection").
				SetTitle("Connection String").
				SetDescription("Set Crate connection string").
				SetMust(true).
				SetDefault("postgresql://crate@localhost:5432/doc?sslmode=disable"),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_idle_connections").
				SetDescription("Set Crate max idle connections").
				SetMust(false).
				SetDefault("10").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("max_open_connections").
				SetDescription("Set Crate max open connections").
				SetMust(false).
				SetDefault("100").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddProperty(
			common.NewProperty().
				SetKind("int").
				SetName("connection_max_lifetime_seconds").
				SetTitle("Connection Lifetime (Seconds)").
				SetDescription("Set Crate connection max lifetime seconds").
				SetMust(false).
				SetDefault("3600").
				SetMin(1).
				SetMax(math.MaxInt32),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("method").
				SetKind("string").
				SetDescription("Set Crate execution method").
				SetOptions([]string{"query", "exec"}).
				SetDefault("query").
				SetMust(true),
		).
		AddMetadata(
			common.NewMetadata().
				SetName("isolation_level").
				SetKind("string").
				SetDescription("Set Crate isolation level").
				SetOptions([]string{"Default", "ReadUncommitted", "ReadCommitted", "RepeatableRead", "Serializable"}).
				SetDefault("Default").
				SetMust(false),
		)

}
