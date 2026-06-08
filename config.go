package main

var tailnet = &tailnetConfig{
	DNS: &dnsConfig{
		MagicDNS:    boolPtr(true),
		Nameservers: []nameserverConfig{},
		SearchPaths: []string{},
		SplitDNS:    []splitDNSConfig{},
	},
	Settings: &tailnetSettingsConfig{
		AclsExternallyManagedOn:               boolPtr(false),
		DevicesApprovalOn:                     boolPtr(false),
		DevicesAutoUpdatesOn:                  boolPtr(true),
		DevicesKeyDurationDays:                intPtr(180),
		HTTPSEnabled:                          boolPtr(false),
		NetworkFlowLoggingOn:                  boolPtr(false),
		PostureIdentityCollectionOn:           boolPtr(false),
		RegionalRoutingOn:                     boolPtr(false),
		UsersApprovalOn:                       boolPtr(true),
		UsersRoleAllowedToJoinExternalTailnet: stringPtr("admin"),
	},
}

var devices = []deviceConfig{
	{
		Name:              "apple-tv",
		DeviceID:          "nF8qDrXfb921CNTRL",
		Hostname:          "apple-tv",
		FQDN:              "apple-tv.tailb35748.ts.net",
		Tags:              []string{"tag:untrusted"},
		KeyExpiryDisabled: false,
		Authorized:        true,
	},
	{
		Name:              "blink",
		DeviceID:          "n28jZ9kBFw11CNTRL",
		Hostname:          "blink",
		FQDN:              "blink.tailb35748.ts.net",
		Tags:              []string{"tag:server"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
	{
		Name:              "bn-mac-edi-008",
		DeviceID:          "nL5mP7qQsd11CNTRL",
		Hostname:          "BN-MAC-EDI-008",
		FQDN:              "bn-mac-edi-008.tailb35748.ts.net",
		Tags:              []string{"tag:personal"},
		KeyExpiryDisabled: false,
		Authorized:        true,
	},
	{
		Name:              "ephone",
		DeviceID:          "ng1XVChP1b11CNTRL",
		Hostname:          "localhost",
		FQDN:              "ephone.tailb35748.ts.net",
		Tags:              []string{"tag:personal"},
		KeyExpiryDisabled: false,
		Authorized:        true,
	},
	{
		Name:              "falcon",
		DeviceID:          "nq5gdWQR7X11CNTRL",
		Hostname:          "falcon",
		FQDN:              "falcon.tailb35748.ts.net",
		Tags:              []string{"tag:server"},
		SubnetRoutes:      []string{"0.0.0.0/0", "::/0"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
	{
		Name:              "fourth",
		DeviceID:          "nexsMcMda321CNTRL",
		Hostname:          "fourth",
		FQDN:              "fourth.tailb35748.ts.net",
		Tags:              []string{"tag:server"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
	{
		Name:              "mba",
		DeviceID:          "nk9XCktEzW11CNTRL",
		Hostname:          "Edward's MacBook Air (2)",
		FQDN:              "mba.tailb35748.ts.net",
		Tags:              []string{"tag:personal"},
		KeyExpiryDisabled: false,
		Authorized:        true,
	},
	{
		Name:              "partridge",
		DeviceID:          "nKBdngosQB21CNTRL",
		Hostname:          "partridge",
		FQDN:              "partridge.tailb35748.ts.net",
		Tags:              []string{"tag:server", "tag:ci-allowed"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
}

func boolPtr(value bool) *bool {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func stringPtr(value string) *string {
	return &value
}
