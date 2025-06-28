package lakekeeper

import "context"

type ServerInfo struct {
	AuthzBackend                 string   `json:"authz-backend"`
	Bootstrapped                 bool     `json:"bootstrapped"`
	DefaultProjectID             string   `json:"default-project-id"`
	AWSSystemIdentitiesEnabled   bool     `json:"aws-system-identities-enabled"`
	AzureSystemIdentitiesEnabled bool     `json:"azure-system-identities-enabled"`
	GCPSystemIdentitiesEnabled   bool     `json:"gcp-system-identities-enabled"`
	ServerID                     string   `json:"server-id"`
	Version                      string   `json:"version"`
	Queues                       []string `json:"queues"`
}

func (lakekeeperClient *Client) GetServerInfo(ctx context.Context) (*ServerInfo, error) {
	var serverInfo ServerInfo
	err := lakekeeperClient.get(ctx, "/management/v1/info", &serverInfo, nil)
	if err != nil {
		return nil, err
	}

	return &serverInfo, nil
}
