# tailnet

Tailscale tailnet management with Pulumi.

## Local setup

Install Go and Pulumi locally, then fetch the Go modules:

```sh
go mod download
```

Create/select the stack:

```sh
pulumi stack init main
pulumi config set tailscale:tailnet <tailnet-name>
```

For local previews and applies, provide credentials through the environment:

```sh
export TAILSCALE_OAUTH_CLIENT_ID=<oauth-client-id>
export TAILSCALE_OAUTH_CLIENT_SECRET=<oauth-client-secret>
pulumi preview
```

OAuth is preferred over a long-lived API key. The Pulumi Tailscale provider also supports `TAILSCALE_API_KEY` if we need to use an API key for the initial import.

## Import existing state

Replace `policy.hujson` and the values in `config.go` with the current values from Tailscale before importing. The policy resource overwrites the whole Tailscale policy file, and `dns` manages the whole DNS configuration.

Import the tailnet-wide resources:

```sh
pulumi import tailscale:index/acl:Acl policy acl
pulumi import tailscale:index/dnsConfiguration:DnsConfiguration dns dns_configuration
pulumi import tailscale:index/tailnetSettings:TailnetSettings settings tailnet_settings
```

Import managed device resources using each device's node ID:

```sh
pulumi import tailscale:index/deviceKey:DeviceKey device-key-<name> <node-id>
pulumi import tailscale:index/deviceTags:DeviceTags device-tags-<name> <node-id>
pulumi import tailscale:index/deviceAuthorization:DeviceAuthorization device-authorization-<name> <node-id>
pulumi import tailscale:index/deviceSubnetRoutes:DeviceSubnetRoutes device-routes-<name> <node-id>
```

Then verify:

```sh
go test ./...
pulumi preview
```

## GitHub Actions

The workflow previews Pulumi changes on pull requests and applies them on pushes to `main`.

Repository secrets:

- `PULUMI_ACCESS_TOKEN`
- `TAILSCALE_OAUTH_CLIENT_ID`
- `TAILSCALE_OAUTH_CLIENT_SECRET`

Repository variables:

- `PULUMI_STACK`, for example `main` or `<pulumi-org>/tailnet/main`
