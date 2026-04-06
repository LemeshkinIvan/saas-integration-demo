package constants

type AmoConfig struct {
	ChannelID           string `yaml:"channel_id"`
	ChannelIDOauth      string `yaml:"channel_id_oauth"`
	ChannelSecret       string `yaml:"channel_secret"`
	ChannelSecretKey    string `yaml:"channel_secret_key"`
	RedirectURL         string `yaml:"redirect_url"`
	ApiConnectMethod    string `yaml:"api_connect_method"`
	ApiDisconnectMethod string `yaml:"api_disconnect_method"`
	Domain              string `yaml:"domain"`
	Referer             string `yaml:"referer"`
	SignatureSecret     string `yaml:"signature_secret"`
	SourceURL           string `yaml:"source_url"`
	PipelineURL         string `yaml:"pipeline_url"`
}
