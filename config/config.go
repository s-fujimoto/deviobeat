// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Deviobeat DeviobeatConfig
}

type DeviobeatConfig struct {
	Period string `yaml:"period"`
	URL string `yaml:"url"`
}
