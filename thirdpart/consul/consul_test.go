package consul

import (
	"fmt"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
)

type (
	// Config -.
	Config struct {
		App    `yaml:"app"`
		Log    `yaml:"logger"`
		Consul `yaml:"consul"`
		PG     `yaml:"postgres"`
		RMQ    `yaml:"rabbitmq"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// Consul -.
	Consul struct {
		CheckApi string `env-required:"true" yaml:"checkApi"    env:"CONSUL_CHECKAPI"`
		Interval string `env-required:"true" yaml:"interval" env:"CONSUL_INTERVAL"`
		Timeout  string `env-required:"true" yaml:"timeout" env:"CONSUL_TIMEOUT"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url"      env:"PG_URL"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true" yaml:"url"                 env:"RMQ_URL"`
	}
)

func Test_Consul(t *testing.T) {
	type args struct {
		consulAddr  string
		serviceName string
		addr        string
		folder      string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Test_Consul", args{"172.16.12.11:8500", "clean", ":8888", "dev"}},
	}
	for _, tt := range tests {
		cfg := &Config{}
		t.Run(tt.name, func(t *testing.T) {
			consulClient, serviceID, _, err := RegisterAndCfgConsul(cfg, tt.args.consulAddr, tt.args.serviceName,
				tt.args.addr, tt.args.folder)
			if err != nil {
				t.Errorf("RegisterAndCfgConsul returns error: %v.", err)
			} else {
				fmt.Println(serviceID)
			}
			DeregisterService(consulClient, serviceID)
		})
	}
}

func Test_ConsulKV(t *testing.T) {
	type kvArgs struct {
		cfg          *Config
		consulClient *consulapi.Client
		folder       string
		serviceName  string
	}
	consulClient, _ := CreateClient("172.16.12.11:8500")
	tests := []struct {
		name string
		args kvArgs
	}{
		{"test_ConsulKV", kvArgs{&Config{}, consulClient, "dev", "user"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consulOption, err := GetKV(tt.args.cfg, tt.args.consulClient, tt.args.folder, tt.args.serviceName)
			if err != nil {
				t.Errorf("GetKV() returns error: %v", err)
			} else {
				t.Logf("%v.", consulOption)
			}
		})
	}
}

func Test_RetrieveAddressPort(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{"test1", args{"HTTP GET http://192.168.2.180:9898/app/healthcheck: 200 OK Output: {\"status\": true}"}, "192.168.2.180", "9898"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RetrieveAddressPort(tt.args.url)
			if got != tt.want {
				t.Errorf("RetrieveAddressPort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RetrieveAddressPort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_Deregister(t *testing.T) {
	type Args struct {
		consulAddr    string
		serviceIDList []string
	}
	tests := []struct {
		name string
		args Args
	}{
		{
			"Test_Deregister",
			Args{
				"172.16.12.11:8500",
				[]string{
					"clean_172.16.12.8:8811",
					"clean_172.16.12.8:8812",
					"clean_172.16.12.8:8813",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consulClient, err := CreateClient(tt.args.consulAddr)
			if err != nil {
				t.Errorf("fail to connect consul, error: %s.", err.Error())
			} else {
				for _, serviceID := range tt.args.serviceIDList {
					err = DeregisterService(consulClient, serviceID)
					if err != nil {
						t.Errorf("fail to deregister service(%s), error: %s.", serviceID, err)
					}
				}
			}
		})
	}
}

func TestRetrieveServiceName(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"testRetrieveServiceName", args{"Service 'helloapp' check"}, "helloapp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveServiceName(tt.args.url); got != tt.want {
				t.Errorf("RetrieveServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveServiceID(t *testing.T) {
	type args struct {
		checkID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"testRetrieveServiceID", args{"service:helloapp-3c17f434b5c34311a068d710624c308a"}, "helloapp-3c17f434b5c34311a068d710624c308a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveServiceID(tt.args.checkID); got != tt.want {
				t.Errorf("RetrieveServiceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
