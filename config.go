package main

type JobConfig struct {
	Name    string `yaml:"name"`
	RunAt   string `yaml:"run_at"`
	Days    []int  `yaml:"days"`
	Command string `yaml:"command"`
}

type ConfigRoot struct {
	Jobs []JobConfig `yaml:"jobs"`
}
