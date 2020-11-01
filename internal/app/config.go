package app

type Config struct {
	Generators  []GeneratorConfig
	Queue       QueueConfig
	Aggregators []AggregatorConfig
	Storage     StorageConfig
}

type GeneratorConfig struct {
	TimeoutSec    int
	SendPeriodSec int
	DataSources   []DataSourceConfig
}

type DataSourceConfig struct {
	Id            string
	InitValue     int
	MaxChangeStep int
}

type QueueConfig struct {
	Size int
}

type AggregatorConfig struct {
	AggregationPeriodSec int
	SubIds               []string
}

type StorageConfig struct {
	Type    int
	Options *StorageOptions
}

type StorageOptions struct {
	Filename string
}
