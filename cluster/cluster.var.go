package cluster

func (client *ClusterClient) GetSourceConfig(typeName string, name string) (config string, err error) {
	dataMap := client.dataMap.Copy()
	dataMap.Set("type", typeName)
	dataMap.Set("name", name)
	values, err := client.handler.GetValue(dataMap.Translate(p_varConfig))
	if err != nil {
		return
	}
	config = string(values)
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