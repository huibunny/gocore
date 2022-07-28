package consul

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v2"
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
func ServicesByPassing(client *consulapi.Client, serviceName string) []PassingService {

	healthChecks, _, _ := client.Health().State(consulapi.HealthPassing, nil)
	var result []PassingService
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
		}
	}

	return result
}

// LogServiceID Log service id
func LogServiceID(services []PassingService) {
	for _, service := range services {
		print("service", service.ServiceID)
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

func RegisterAndCfgConsul(cfg interface{}, consulAddr, serviceName,
	host, port, consulInterval,
	consulTimeout, folder string) (*consulapi.Client, string, error) {
	// 创建consul api客户端
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		os.Exit(1)
	}

	var serviceID string
	serviceID, err = RegisterService(serviceName, *consulClient, host, port, consulInterval, consulTimeout)
	if err == nil {
		key := strings.Join([]string{folder, serviceName}, "/")
		kv, _, err := consulClient.KV().Get(key, nil)
		if err == nil {
			// only support yaml kv
			err = yaml.NewDecoder(strings.NewReader(string(kv.Value))).Decode(cfg)
			if err == nil {
			} else {
				print("error: " + err.Error())
			}
		} else {
			print("error: " + err.Error())
		}
	} else {
		print("error: " + err.Error())
	}
	return consulClient, serviceID, err
}

// RegisterService register service in consul
func RegisterService(service string, client consulapi.Client,
	svcHost string, svcPort string, consulInterval string,
	consulTimeout string) (string, error) {
	svcAddress := svcHost + ":" + svcPort

	// 设置Consul对服务健康检查的参数
	check := consulapi.AgentServiceCheck{
		HTTP:     "http://" + svcAddress + "/healthz",
		Interval: consulInterval + "s",
		Timeout:  consulTimeout + "s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	//设置微服务Consul的注册信息
	reg := &consulapi.AgentServiceRegistration{
		ID:      service + "_" + svcAddress,
		Name:    service,
		Address: svcHost,
		Port:    port,
		Check:   &check,
	}

	// 执行注册
	err := client.Agent().ServiceRegister(reg)

	return reg.ID, err
}

func DeregisterService(consulClient *consulapi.Client, serviceID string) {
	consulClient.Agent().ServiceDeregister(serviceID)
}
