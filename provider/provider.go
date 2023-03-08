package provider

import (
	"fmt"

	"github.com/DeltaLaboratory/dynamicdns/dnsapi"
	"github.com/DeltaLaboratory/dynamicdns/provider/cloudflare"
)

type Provider struct {
	providers map[string]dnsapi.APICreator
}

func NewProvider() *Provider {
	p := new(Provider)
	p.providers = map[string]dnsapi.APICreator{}
	p.SetProvider("cloudflare", cloudflare.NewAPI)
	return p
}

func (p *Provider) SetProvider(name string, apiCreator dnsapi.APICreator) {
	p.providers[name] = apiCreator
}

func (p *Provider) GetProvider(name string) (dnsapi.APICreator, error) {
	v, ok := p.providers[name]
	if !ok {
		return nil, fmt.Errorf("no such provider")
	}
	return v, nil
}
