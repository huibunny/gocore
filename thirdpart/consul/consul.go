package consul

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/huibunny/gocore/utils"
	"github.com/spf13/viper"
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

func CreateClient(consulAddr string) *consulapi.Client {
	// 创建consul api客户端
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulAddr
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		os.Exit(1)
	}

	return consulClient
}

func GetKV(cfg interface{}, consulClient *consulapi.Client, folder, serviceName string) (map[string]string, error) {
	var consulOption map[string]string
	key := strings.Join([]string{folder, serviceName}, "/")
	kv, _, err := consulClient.KV().Get(key, nil)
	if err == nil {
		if kv == nil {
			err = errors.New("KV not found for " + key + ".")
		} else {
			// only support yaml kv
			kvIO := strings.NewReader(string(kv.Value))
			err = yaml.NewDecoder(kvIO).Decode(cfg)
			viper.SetConfigType("yaml")
			err := viper.ReadConfig(bytes.NewBuffer(kv.Value))
			if err != nil {
				print(err)
			} else {
				consulOption = viper.GetStringMapString("consul")
			}
		}
	} else {
	}

	return consulOption, err
}

func RegisterAndCfgConsul(cfg interface{}, consulAddr, serviceName,
	listenAddr, folder string) (*consulapi.Client, string, string, error) {
	host, port := utils.GetHostPort(listenAddr)
	consulClient := CreateClient(consulAddr)
	var err error
	var serviceID string
	if consulClient == nil {
		var consulOption map[string]string
		consulOption, err = GetKV(cfg, consulClient, folder, serviceName)
		if err == nil {
			serviceID, err = RegisterService(serviceName, *consulClient, host, port, consulOption)
		} else {
			err = errors.New("fail to get kv(" + folder + "/" + serviceName + ") from consul, error: " + err.Error() + ".")
		}
	} else {
		err = errors.New("fail to connect consul(" + consulAddr + ").")
	}

	return consulClient, serviceID, port, err
}

// RegisterService register service in consul
func RegisterService(service string, client consulapi.Client,
	host, port string, consulOption map[string]string) (string, error) {
	svcAddress := strings.Join([]string{host, port}, ":")

	checkApi := consulOption["checkapi"]
	interval := consulOption["interval"]
	timeout := consulOption["timeout"]
	// 设置Consul对服务健康检查的参数
	if strings.HasPrefix(checkApi, "/") {
		//
	} else {
		checkApi = strings.Join([]string{"/", checkApi}, "")
	}
	check := consulapi.AgentServiceCheck{
		HTTP:     "http://" + svcAddress + checkApi,
		Interval: interval + "s",
		Timeout:  timeout + "s",
		Notes:    "Consul check service health status.",
	}

	intPort, _ := strconv.Atoi(port)

	//设置微服务Consul的注册信息
	serviceID := service + "_" + svcAddress
	reg := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    service,
		Address: host,
		Port:    intPort,
		Check:   &check,
	}

	// 执行注册
	err := client.Agent().ServiceRegister(reg)

	return serviceID, err
}

func DeregisterService(consulClient *consulapi.Client, serviceID string) {
	consulClient.Agent().ServiceDeregister(serviceID)
}
