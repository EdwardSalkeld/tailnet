package main

import (
	"os"

	"github.com/pulumi/pulumi-tailscale/sdk/go/tailscale"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type nameserverConfig struct {
	Address         string `json:"address"`
	UseWithExitNode *bool  `json:"useWithExitNode,omitempty"`
}

type splitDNSConfig struct {
	Domain      string             `json:"domain"`
	Nameservers []nameserverConfig `json:"nameservers"`
}

type dnsConfig struct {
	MagicDNS         *bool              `json:"magicDns,omitempty"`
	OverrideLocalDNS *bool              `json:"overrideLocalDns,omitempty"`
	Nameservers      []nameserverConfig `json:"nameservers,omitempty"`
	SearchPaths      []string           `json:"searchPaths,omitempty"`
	SplitDNS         []splitDNSConfig   `json:"splitDns,omitempty"`
}

type tailnetSettingsConfig struct {
	AclsExternalLink                      *string `json:"aclsExternalLink,omitempty"`
	AclsExternallyManagedOn               *bool   `json:"aclsExternallyManagedOn,omitempty"`
	DevicesApprovalOn                     *bool   `json:"devicesApprovalOn,omitempty"`
	DevicesAutoUpdatesOn                  *bool   `json:"devicesAutoUpdatesOn,omitempty"`
	DevicesKeyDurationDays                *int    `json:"devicesKeyDurationDays,omitempty"`
	HTTPSEnabled                          *bool   `json:"httpsEnabled,omitempty"`
	NetworkFlowLoggingOn                  *bool   `json:"networkFlowLoggingOn,omitempty"`
	PostureIdentityCollectionOn           *bool   `json:"postureIdentityCollectionOn,omitempty"`
	RegionalRoutingOn                     *bool   `json:"regionalRoutingOn,omitempty"`
	UsersApprovalOn                       *bool   `json:"usersApprovalOn,omitempty"`
	UsersRoleAllowedToJoinExternalTailnet *string `json:"usersRoleAllowedToJoinExternalTailnet,omitempty"`
}

type tailnetConfig struct {
	DNS      *dnsConfig             `json:"dns,omitempty"`
	Settings *tailnetSettingsConfig `json:"settings,omitempty"`
}

type deviceConfig struct {
	Name              string   `json:"name"`
	DeviceID          string   `json:"deviceId"`
	Hostname          string   `json:"hostname,omitempty"`
	FQDN              string   `json:"fqdn,omitempty"`
	Tags              []string `json:"tags"`
	SubnetRoutes      []string `json:"subnetRoutes,omitempty"`
	KeyExpiryDisabled bool     `json:"keyExpiryDisabled"`
	Authorized        bool     `json:"authorized"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		var policyResource *tailscale.Acl

		policy, err := os.ReadFile("policy.hujson")
		if err == nil {
			policyResource, err = tailscale.NewAcl(ctx, "policy", &tailscale.AclArgs{
				Acl:               pulumi.String(string(policy)),
				ResetAclOnDestroy: pulumi.Bool(false),
			}, pulumi.Protect(true))
			if err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		if tailnet.DNS != nil {
			_, err = tailscale.NewDnsConfiguration(ctx, "dns", buildDNSArgs(tailnet.DNS), pulumi.Protect(true))
			if err != nil {
				return err
			}
		}

		if tailnet.Settings != nil {
			_, err = tailscale.NewTailnetSettings(ctx, "settings", buildTailnetSettingsArgs(tailnet.Settings), pulumi.Protect(true))
			if err != nil {
				return err
			}
		}

		for _, device := range devices {
			deviceTagOptions := []pulumi.ResourceOption{pulumi.Protect(true)}
			if policyResource != nil {
				deviceTagOptions = append(deviceTagOptions, pulumi.DependsOn([]pulumi.Resource{policyResource}))
			}

			_, err = tailscale.NewDeviceKey(ctx, "device-key-"+device.Name, &tailscale.DeviceKeyArgs{
				DeviceId:          pulumi.String(device.DeviceID),
				KeyExpiryDisabled: pulumi.Bool(device.KeyExpiryDisabled),
			}, pulumi.Protect(true))
			if err != nil {
				return err
			}

			_, err = tailscale.NewDeviceTags(ctx, "device-tags-"+device.Name, &tailscale.DeviceTagsArgs{
				DeviceId: pulumi.String(device.DeviceID),
				Tags:     toStringArray(device.Tags),
			}, deviceTagOptions...)
			if err != nil {
				return err
			}

			if device.SubnetRoutes != nil {
				_, err = tailscale.NewDeviceSubnetRoutes(ctx, "device-routes-"+device.Name, &tailscale.DeviceSubnetRoutesArgs{
					DeviceId: pulumi.String(device.DeviceID),
					Routes:   toStringArray(device.SubnetRoutes),
				}, pulumi.Protect(true))
				if err != nil {
					return err
				}
			}

			_, err = tailscale.NewDeviceAuthorization(ctx, "device-authorization-"+device.Name, &tailscale.DeviceAuthorizationArgs{
				DeviceId:   pulumi.String(device.DeviceID),
				Authorized: pulumi.Bool(device.Authorized),
			}, pulumi.Protect(true))
			if err != nil {
				return err
			}
		}

		ctx.Export("managedResources", pulumi.Map{
			"policy":               pulumi.Bool(policy != nil),
			"dns":                  pulumi.Bool(tailnet.DNS != nil),
			"settings":             pulumi.Bool(tailnet.Settings != nil),
			"deviceKeys":           pulumi.Int(len(devices)),
			"deviceTags":           pulumi.Int(len(devices)),
			"deviceSubnetRoutes":   pulumi.Int(countDeviceSubnetRoutes(devices)),
			"deviceAuthorizations": pulumi.Int(len(devices)),
		})

		return nil
	})
}

func countDeviceSubnetRoutes(devices []deviceConfig) int {
	count := 0
	for _, device := range devices {
		if device.SubnetRoutes != nil {
			count++
		}
	}
	return count
}

func buildDNSArgs(config *dnsConfig) *tailscale.DnsConfigurationArgs {
	args := &tailscale.DnsConfigurationArgs{}

	if config.MagicDNS != nil {
		args.MagicDns = pulumi.Bool(*config.MagicDNS)
	}
	if config.OverrideLocalDNS != nil {
		args.OverrideLocalDns = pulumi.Bool(*config.OverrideLocalDNS)
	}
	if config.Nameservers != nil {
		args.Nameservers = buildNameservers(config.Nameservers)
	}
	if config.SearchPaths != nil {
		args.SearchPaths = toStringArray(config.SearchPaths)
	}
	if config.SplitDNS != nil {
		args.SplitDns = buildSplitDNS(config.SplitDNS)
	}

	return args
}

func buildNameservers(configs []nameserverConfig) tailscale.DnsConfigurationNameserverArray {
	result := make(tailscale.DnsConfigurationNameserverArray, 0, len(configs))
	for _, config := range configs {
		args := &tailscale.DnsConfigurationNameserverArgs{
			Address: pulumi.String(config.Address),
		}
		if config.UseWithExitNode != nil {
			args.UseWithExitNode = pulumi.Bool(*config.UseWithExitNode)
		}
		result = append(result, args)
	}
	return result
}

func buildSplitDNS(configs []splitDNSConfig) tailscale.DnsConfigurationSplitDnArray {
	result := make(tailscale.DnsConfigurationSplitDnArray, 0, len(configs))
	for _, config := range configs {
		result = append(result, &tailscale.DnsConfigurationSplitDnArgs{
			Domain:      pulumi.String(config.Domain),
			Nameservers: buildSplitDNSNameservers(config.Nameservers),
		})
	}
	return result
}

func buildSplitDNSNameservers(configs []nameserverConfig) tailscale.DnsConfigurationSplitDnNameserverArray {
	result := make(tailscale.DnsConfigurationSplitDnNameserverArray, 0, len(configs))
	for _, config := range configs {
		args := &tailscale.DnsConfigurationSplitDnNameserverArgs{
			Address: pulumi.String(config.Address),
		}
		if config.UseWithExitNode != nil {
			args.UseWithExitNode = pulumi.Bool(*config.UseWithExitNode)
		}
		result = append(result, args)
	}
	return result
}

func buildTailnetSettingsArgs(config *tailnetSettingsConfig) *tailscale.TailnetSettingsArgs {
	args := &tailscale.TailnetSettingsArgs{}

	if config.AclsExternalLink != nil {
		args.AclsExternalLink = pulumi.String(*config.AclsExternalLink)
	}
	if config.AclsExternallyManagedOn != nil {
		args.AclsExternallyManagedOn = pulumi.Bool(*config.AclsExternallyManagedOn)
	}
	if config.DevicesApprovalOn != nil {
		args.DevicesApprovalOn = pulumi.Bool(*config.DevicesApprovalOn)
	}
	if config.DevicesAutoUpdatesOn != nil {
		args.DevicesAutoUpdatesOn = pulumi.Bool(*config.DevicesAutoUpdatesOn)
	}
	if config.DevicesKeyDurationDays != nil {
		args.DevicesKeyDurationDays = pulumi.Int(*config.DevicesKeyDurationDays)
	}
	if config.HTTPSEnabled != nil {
		args.HttpsEnabled = pulumi.Bool(*config.HTTPSEnabled)
	}
	if config.NetworkFlowLoggingOn != nil {
		args.NetworkFlowLoggingOn = pulumi.Bool(*config.NetworkFlowLoggingOn)
	}
	if config.PostureIdentityCollectionOn != nil {
		args.PostureIdentityCollectionOn = pulumi.Bool(*config.PostureIdentityCollectionOn)
	}
	if config.RegionalRoutingOn != nil {
		args.RegionalRoutingOn = pulumi.Bool(*config.RegionalRoutingOn)
	}
	if config.UsersApprovalOn != nil {
		args.UsersApprovalOn = pulumi.Bool(*config.UsersApprovalOn)
	}
	if config.UsersRoleAllowedToJoinExternalTailnet != nil {
		args.UsersRoleAllowedToJoinExternalTailnet = pulumi.String(*config.UsersRoleAllowedToJoinExternalTailnet)
	}

	return args
}

func toStringArray(values []string) pulumi.StringArray {
	result := make(pulumi.StringArray, 0, len(values))
	for _, value := range values {
		result = append(result, pulumi.String(value))
	}
	return result
}
