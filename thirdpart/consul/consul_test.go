package consul

import (
	"testing"
)

func Test_Consul(t *testing.T) {
	type (

		// App -.
		App struct {
			Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
			Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
		}

		// Log -.
		Log struct {
			Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
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

		// Config -.
		Config struct {
			App `yaml:"app"`
			Log `yaml:"logger"`
			PG  `yaml:"postgres"`
			RMQ `yaml:"rabbitmq"`
		}
	)

	type args struct {
		consulAddr     string
		serviceName    string
		host           string
		port           string
		consulInterval string
		consulTimeout  string
		folder         string
	}

	tests := []struct {
		name string
		args args
	}{
		{"Test_Consul", args{"127.0.0.1:8500", "clean", "172.16.12.8", "8888", "3", "3", "dev"}},
	}
	for _, tt := range tests {
		cfg := &Config{}
		t.Run(tt.name, func(t *testing.T) {
			consulClient, serviceID, err := RegisterAndCfgConsul(cfg, tt.args.consulAddr, tt.args.serviceName,
				tt.args.host, tt.args.port, tt.args.consulInterval, tt.args.consulTimeout, tt.args.folder)
			if err != nil {
				t.Errorf("RegisterAndCfgConsul returns error: %v.", err)
			}
			DeregisterService(consulClient, serviceID)
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
