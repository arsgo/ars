package cluster

func (client *ClusterClient) GetSourceConfig(typeName string, name string) (config string, err error) {

	dataMap := client.dataMap.Copy()
	dataMap.Set("type", typeName)
	dataMap.Set("name", name)
	path := dataMap.Translate(p_varConfig)

	cfg := client.configCache.Get(path)
	if cfg != nil {
		config = cfg.(string)
		return
	}
	values, err := client.handler.GetValue(path)
	if err != nil {
		return
	}
	config = string(values)
	client.configCache.Set(path, config)
	client.WatchClusterValueChange(path, func() {
		values, err := client.handler.GetValue(path)
		if err != nil {
			return
		}
		client.configCache.Set(path, string(values))
	})

	return
}

func (client *ClusterClient) GetMQConfig(name string) (string, error) {
	return client.GetSourceConfig("mq", name)
}
func (client *ClusterClient) GetElasticConfig(name string) (string, error) {
	return client.GetSourceConfig("elastic", name)
}
func (client *ClusterClient) GetDBConfig(name string) (string, error) {
	return client.GetSourceConfig("db", name)
}
