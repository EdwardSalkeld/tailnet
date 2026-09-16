package main

var tailnet = &tailnetConfig{
	DNS: &dnsConfig{
		MagicDNS:    boolPtr(true),
		Nameservers: []nameserverConfig{},
		SearchPaths: []string{},
		// The home resolver is reached through Partridge's approved subnet route.
		// It serves the private int.alcachofa.faith zone, while public DNS remains
		// responsible for every other name.
		//
		// UseWithExitNode keeps this resolver in play while an exit node is
		// selected. Without it a client on an exit node sends every query to
		// that node's resolver, and the int zone stops resolving away from home.
		SplitDNS: []splitDNSConfig{
			{
				Domain: "int.alcachofa.faith",
				Nameservers: []nameserverConfig{
					{Address: "10.4.1.1", UseWithExitNode: boolPtr(true)},
				},
			},
		},
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
		DeviceID:          "nCGPKNNscf11CNTRL",
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
		Tags:              []string{"tag:server", "tag:ci-allowed"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
	{
		Name:              "kite",
		DeviceID:          "n1E48WvbR111CNTRL",
		Hostname:          "kite",
		FQDN:              "kite.tailb35748.ts.net",
		Tags:              []string{"tag:server"},
		KeyExpiryDisabled: true,
		Authorized:        true,
	},
	{
		Name:              "magpie",
		DeviceID:          "n49eMmw3F311CNTRL",
		Hostname:          "magpie",
		FQDN:              "magpie.tailb35748.ts.net",
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
		Name:     "partridge",
		DeviceID: "nKBdngosQB21CNTRL",
		Hostname: "partridge",
		FQDN:     "partridge.tailb35748.ts.net",
		Tags:     []string{"tag:server"},
		// The NixOS host advertises the less-specific 10.4.0.0/23 so clients
		// already on 10.4.1.0/24 retain their direct LAN route. Policy below
		// grants access only to the actual 10.4.1.0/24 LAN.
		//
		// 0.0.0.0/0 and ::/0 approve Partridge as an exit node. Home has no
		// IPv6 upstream, so v6 destinations fall back to v4 rather than
		// leaking out of the client's local interface.
		SubnetRoutes:      []string{"10.4.0.0/23", "0.0.0.0/0", "::/0"},
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
