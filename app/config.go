package app

type Config struct {
	Generators  []GeneratorConfig  `json:"generators"`
	Queue       QueueConfig        `json:"queue"`
	Aggregators []AggregatorConfig `json:"aggregators"`
	StorageType int                `json:"storage_type"`
}

type GeneratorConfig struct {
	TimeoutSec    int                `json:"timeout_s"`
	SendPeriodSec int                `json:"send_period_s"`
	DataSources   []DataSourceConfig `json:"data_sources"`
}

type DataSourceConfig struct {
	Id            string `json:"id"`
	InitValue     int    `json:"init_value"`
	MaxChangeStep int    `json:"max_change_step"`
}

type QueueConfig struct {
	Size int
}

type AggregatorConfig struct {
	AggregationPeriodSec int      `json:"aggregation_period_s"`
	SubIds               []string `json:"sub_ids"`
}
