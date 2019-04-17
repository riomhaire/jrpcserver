package consulagent

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/riomhaire/consul"
)

type ConsulServiceRegistry struct {
	name            string
	applicationHost string
	applicationPort int
	consulHost      string

	baseEndpoint   string
	healthEndpoint string
	id             string

	consulClient *consul.ConsulClient // This registers this service with consul - may extract this into a separate use case

}

func NewConsulServiceRegistry(consulHost, name, applicationHost string, applicationPort int, baseEndpoint, healthEndpoint string) *ConsulServiceRegistry {
	r := ConsulServiceRegistry{}
	r.baseEndpoint = baseEndpoint
	r.healthEndpoint = healthEndpoint
	r.name = name
	r.consulHost = consulHost
	r.applicationHost = applicationHost
	r.applicationPort = applicationPort

	return &r
}

func (a *ConsulServiceRegistry) Register() error {
	// Register with consol (if required)
	id := fmt.Sprintf("%v-%v-%v", a.name, a.applicationHost, a.applicationPort)
	a.id = id // This is our safe copy

	a.consulClient, _ = consul.NewConsulClient(a.consulHost)
	health := fmt.Sprintf("http://%v:%v%v", a.applicationHost, a.applicationPort, a.healthEndpoint)
	log.Printf("Registering with Consul at %v with %v %v\n", a.consulHost, a.baseEndpoint, health)

	a.consulClient.PeriodicRegister(id, a.name, a.applicationHost, a.applicationPort, a.baseEndpoint, health, 15)
	return nil

}

/*

 */
func (a *ConsulServiceRegistry) Deregister() error {
	log.Printf("De Registering %v with Consul at %v with %v \n", a.id, a.consulHost, a.baseEndpoint)
	a.consulClient.DeRegister(a.id)
	return nil
}
