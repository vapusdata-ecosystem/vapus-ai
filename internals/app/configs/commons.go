package appconfigs

type VapusSvcNetworkConfig struct {
	Port        int32  `yaml:"port"`
	ServiceName string `yaml:"serviceName"`
	// Scheme      string `yaml:"scheme"`
	// ExternalURL string `yaml:"externalUrl"`
	SvcType string `yaml:"svcType"`
	// Addr        string `yaml:"addr"`
	NodePort    int32 `yaml:"nodePort"`
	ServicePort int32 `yaml:"servicePort"`
	HttpGwPort  int32 `yaml:"httpGwPort"`
}
type NetworkConfig struct {
	ExternalURL string                 `yaml:"externalUrl"`
	AIStudioSvc *VapusSvcNetworkConfig `yaml:"aistudioSvc"`
	WebAppSvc   *VapusSvcNetworkConfig `yaml:"webappSvc"`
	AIGateway   *VapusSvcNetworkConfig `yaml:"aigateway"`
	AIUtility   *VapusSvcNetworkConfig `yaml:"aiutility"`
	GatewayURL  string                 `yaml:"gatewayUrl"`
}

type SecretMap struct {
	FilePath string `yaml:"filePath"`
	Secret   string `yaml:"secret"`
	Version  string `yaml:"version"`
}

type TlsCertConfig struct {
	Mtls           bool   `yaml:"mtls"`
	PlainTls       bool   `yaml:"plainTls"`
	Insecure       bool   `yaml:"insecure"`
	CaCertFile     string `yaml:"caCertFile"`
	ServerCertFile string `yaml:"serverCertFile"`
	ServerKeyFile  string `yaml:"serverKeyFile"`
	ClientCertFile string `yaml:"serverCertFile"`
	ClientKeyFile  string `yaml:"serverKeyFile"`
}

type TrinoDeploymentSpec struct {
	Namespace                      string `json:"namespace" yaml:"namespace"`
	AppSelector                    string `json:"appSelector" yaml:"appSelector"`
	TrinoCordDeployment            string `json:"trinoCordDeployment" yaml:"trinoCordDeployment"`
	TrinoWorkerDeployment          string `json:"trinoWorkerDeployment" yaml:"trinoWorkerDeployment"`
	TrinoCordSvc                   string `json:"trinoCordSvc" yaml:"trinoCordSvc"`
	TrinoWorkerSvc                 string `json:"trinoWorkerSvc" yaml:"trinoWorkerSvc"`
	TrinoCordDeploymentContainer   string `json:"trinoCordDeploymentContainer" yaml:"trinoCordDeploymentContainer"`
	TrinoWorkerDeploymentContainer string `json:"trinoWorkerDeploymentContainer" yaml:"trinoWorkerDeploymentContainer"`
	TrnioWorkerPort                int32  `json:"trnioWorkerPort" yaml:"trnioWorkerPort"`
	TrinoCordPort                  int32  `json:"trinoCordPort" yaml:"trinoCordPort"`
	TrinoCordSvcPort               int32  `json:"trinoCordSvcPort" yaml:"trinoCordSvcPort"`
	TrinoWorkerSvcPort             int32  `json:"trinoWorkerSvcPort" yaml:"trinoWorkerSvcPort"`
	TrinoAppName                   string `json:"trinoAppName" yaml:"trinoAppName"`
	TrinoCatalog                   string `json:"trinoCatalog" yaml:"trinoCatalog"`
}
