package common

type ServiceInfo struct {
	DBTableName     string
	ServiceInstance Service
}

// K - ServiceName
// V - ServiceInfo
var RegisteredServices map[string]ServiceInfo = make(map[string]ServiceInfo)
var MainERP ERPSystem

func GetTableNames() []string {
	if len(RegisteredServices) == 0 {
		return nil
	}
	tableNames := make([]string, len(RegisteredServices))
	for _, v := range RegisteredServices {
		tableNames = append(tableNames, v.DBTableName)
	}
	return tableNames
}
