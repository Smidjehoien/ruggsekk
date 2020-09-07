package project

type ProjectConfig struct {
	Id      string `validate:"required"`
	Name    string `validate:"required"`
	Version string `validate:"required"`

	HelmVersion      string `ini:"helm_version"`
	K3sVersion       string `ini:"k3s_version"`
	ServerBinaryName string `ini:"server_binary_name"`

	DockerfileServices []*DockerfileServiceConfig
	HelmServices       []*HelmServiceConfig
	ManifestServices   []*ManifestServiceConfig
}

type DockerfileServiceConfig struct {
	//Common fields
	Id   string `validate:"required"`
	Type string `validate:"required,oneof=dockerfile helm manifest"`
	Port int    `validate:"required"`

	//Unique Fields
	BaseDir         string `validate:"required"`
	Dockerfile      string `validate:"required"`
	ServiceVersion  string `ini:"service_version"`
	UrlPath         string `ini:"base_url"`
	PathPrefixStrip bool   `ini:"path_prefix_strip"`
}

type HelmServiceConfig struct {
	//Common fields
	Id   string `validate:"required"`
	Type string `validate:"required,oneof=dockerfile helm manifest"`
	Port int    `validate:"required"`

	//Unique Fields
	Chart   string `validate:"required"`
	Version string `validate:"required"`
}

type ManifestServiceConfig struct {
	//Common fields
	Id   string `validate:"required"`
	Type string `validate:"required,oneof=dockerfile helm manifest"`
	Port int    `validate:"required"`

	//Unique Fields
	BaseDir  string `validate:"required" ini:"base_dir"`
	Manifest string `validate:"required"`
}
