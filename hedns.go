package hedns

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/UG40A/hedns"
)

// Provider wraps the provider implementation as a Caddy module.
func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.hedns",
		New: func() caddy.Module { return &Provider{new(duckdns.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.Domain = repl.ReplaceAll(p.Provider.Domain, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// duckdns [<api_token>] {
//     api_token <api_token>
//     override_domain <duckdns_domain>
// }
//
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.APIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				p.Provider.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "domain":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.Domain != "" {
					return d.Err("Override domain already set")
				}
				p.Provider.Domain = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
