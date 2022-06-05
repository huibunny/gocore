package consul

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// RetrieveAddressPort 从url中获取地址和端口
// "HTTP GET http://192.168.2.180:9898/app/healthcheck: 200 OK Output: {\"status\": true}"
// "192.168.2.180", "9898"
func RetrieveAddressPort(url string) (string, string) {
	var values = strings.Split(url, "/")
	if len(values) > 2 {
		var addressport = strings.Split(values[2], ":")
		if len(addressport) > 1 {
			return addressport[0], addressport[1]
		}
	}

	return "", ""
}

// RetrieveServiceName 从name中获取Service Name
// "Service 'helloapp' check"
// "helloapp"
func RetrieveServiceName(name string) string {
	var values = strings.Split(name, "'")
	if len(values) > 2 {
		return values[1]
	}

	return ""
}

// RetrieveServiceID 从CheckID中获取ServiceID
// "service:helloapp-3c17f434b5c34311a068d710624c308a"
// "helloapp-3c17f434b5c34311a068d710624c308a"
func RetrieveServiceID(checkID string) string {
	var values = strings.Split(checkID, ":")
	if len(values) > 1 {
		return values[1]
	}

	return ""
}

// PassingService Passing Service structure
type PassingService struct {
	address   string
	port      string
	ServiceID string
}

func ServiceAddress(passingService PassingService) string {
	return fmt.Sprintf("%s:%s", passingService.address, passingService.port)
}

// ServicesByPassing 获取可用的服务
func ServicesByPassing(client *api.Client, serviceName string, logger log.Logger) []PassingService {

	healthChecks, _, _ := client.Health().State(api.HealthPassing, nil)
	var result []PassingService
	if healthChecks != nil {
		logger.Log("len", len(healthChecks))
		for _, healthCheck := range healthChecks {
			name := RetrieveServiceName(healthCheck.Name)
			if name == serviceName {
				address, port := RetrieveAddressPort(healthCheck.Output)
				serviceID := RetrieveServiceID(healthCheck.CheckID)
				service := PassingService{
					address,
					port,
					serviceID,
				}
				result = append(result, service)
				logger.Log("HealthCheck ServiceID", serviceID)
			} else {
				logger.Log("HealthCheck Name", healthCheck.Name)
				logger.Log("servcieName", serviceName)
			}
		}
	}

	return result
}

// LogServiceID Log service id
func LogServiceID(services []PassingService, logger log.Logger) {
	logger.Log("service count", len(services))
	for _, service := range services {
		logger.Log("service", service.ServiceID)
	}
}

// ServiceIndex find service index
func ServiceIndex(pathArray []string, item string) int {
	index := -1

	for i := range pathArray {
		if pathArray[i] == item {
			index = i + 1
			break
		}
	}

	return index
}

// RegisterService register service in consul
func RegisterService(service string, client consul.Client, svcHost string,
	svcPort string, tags string, interval string, logger log.Logger) (registar sd.Registrar) {
	svcAddress := svcHost + ":" + svcPort

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcAddress + "/app/health",
		Interval: interval,
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	//设置微服务Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      service + "_" + svcAddress,
		Name:    service,
		Address: svcHost,
		Port:    port,
		Tags:    []string{service, tags},
		Check:   &check,
	}

	// 执行注册
	registar = consul.NewRegistrar(client, &reg, logger)
	//启动前执行注册
	registar.Register()
	return
}
