// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package health

import (
	"fmt"
)

/**
This file contains definitions for the Warnables maintained within this `health` package.
*/

// updateAvailableWarnable is a Warnable that warns the user that an update is available.
var updateAvailableWarnable = Register(&Warnable{
	Code:     "update-available",
	Title:    "Update available",
	Severity: SeverityLow,
	Text: func(args Args) string {
		return fmt.Sprintf("An update from version %s to %s is available. Run `tailscale update` or `tailscale set --auto-update` to update.", args[ArgCurrentVersion], args[ArgAvailableVersion])
	},
})

// securityUpdateAvailableWarnable is a Warnable that warns the user that an important security update is available.
var securityUpdateAvailableWarnable = Register(&Warnable{
	Code:     "security-update-available",
	Title:    "Security update available",
	Severity: SeverityHigh,
	Text: func(args Args) string {
		return fmt.Sprintf("An urgent security update from version %s to %s is available. Run `tailscale update` or `tailscale set --auto-update` to update now.", args[ArgCurrentVersion], args[ArgAvailableVersion])
	},
})

// unstableWarnable is a Warnable that warns the user that they are using an unstable version of Tailscale
// so they won't be surprised by all the issues that may arise.
var unstableWarnable = Register(&Warnable{
	Code:     "is-using-unstable-version",
	Title:    "Using an unstable version",
	Severity: SeverityLow,
	Text:     StaticMessage("This is an unstable version of Tailscale meant for testing and development purposes: please report any bugs to Tailscale."),
})

// NetworkStatusWarnable is a Warnable that warns the user that the network is down.
var NetworkStatusWarnable = Register(&Warnable{
	Code:                "network-status",
	Title:               "Network down",
	Severity:            SeverityHigh,
	Text:                StaticMessage("Tailscale cannot connect because the network is down. (No network interface is up.)"),
	ImpactsConnectivity: true,
})

// IPNStateWarnable is a Warnable that warns the user that Tailscale is stopped.
var IPNStateWarnable = Register(&Warnable{
	Code:     "wantrunning-false",
	Title:    "Not connected to Tailscale",
	Severity: SeverityLow,
	Text:     StaticMessage("Tailscale is stopped."),
})

// localLogWarnable is a Warnable that warns the user that the local log is misconfigured.
var localLogWarnable = Register(&Warnable{
	Code:     "local-log-config-error",
	Title:    "Local log misconfiguration",
	Severity: SeverityLow,
	Text: func(args Args) string {
		return fmt.Sprintf("The local log is misconfigured: %v", args[ArgError])
	},
})

// LoginStateWarnable is a Warnable that warns the user that they are logged out,
// and provides the last login error if available.
var LoginStateWarnable = Register(&Warnable{
	Code:     "login-state",
	Title:    "Logged out",
	Severity: SeverityMedium,
	Text: func(args Args) string {
		if args[ArgError] != "" {
			return fmt.Sprintf("You are logged out. The last login error was: %v", args[ArgError])
		} else {
			return "You are logged out."
		}
	},
})

// notInMapPollWarnable is a Warnable that warns the user that they cannot connect to the control server.
var notInMapPollWarnable = Register(&Warnable{
	Code:      "not-in-map-poll",
	Title:     "Cannot connect to control server",
	Severity:  SeverityMedium,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text:      StaticMessage("Cannot connect to the control server (not in map poll). Check your Internet connection."),
})

// noDERPHomeWarnable is a Warnable that warns the user that Tailscale doesn't have a home DERP.
var noDERPHomeWarnable = Register(&Warnable{
	Code:      "no-derp-home",
	Title:     "No home relay server",
	Severity:  SeverityHigh,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text:      StaticMessage("Tailscale could not connect to any relay server. Check your Internet connection."),
})

// noDERPConnectionWarnable is a Warnable that warns the user that Tailscale couldn't connect to a specific DERP server.
var noDERPConnectionWarnable = Register(&Warnable{
	Code:      "no-derp-connection",
	Title:     "Relay server unavailable",
	Severity:  SeverityHigh,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text: func(args Args) string {
		if n := args[ArgDERPRegionName]; n != "" {
			return fmt.Sprintf("Tailscale could not connect to the '%s' relay server. Your Internet connection might be down, or the server might be temporarily unavailable.", n)
		} else {
			return fmt.Sprintf("Tailscale could not connect to the relay server with ID '%s'. Your Internet connection might be down, or the server might be temporarily unavailable.", args[ArgDERPRegionID])
		}
	},
})

// derpTimeoutWarnable is a Warnable that warns the user that Tailscale hasn't heard from the home DERP region for a while.
var derpTimeoutWarnable = Register(&Warnable{
	Code:      "derp-timed-out",
	Title:     "Relay server timed out",
	Severity:  SeverityMedium,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text: func(args Args) string {
		if n := args[ArgDERPRegionName]; n != "" {
			return fmt.Sprintf("Tailscale hasn't heard from the '%s' relay server in %v. The server might be temporarily unavailable, or your Internet connection might be down.", n, args[ArgDuration])
		} else {
			return fmt.Sprintf("Tailscale hasn't heard from the home relay server (region ID '%v') in %v. The server might be temporarily unavailable, or your Internet connection might be down.", args[ArgDERPRegionID], args[ArgDuration])
		}
	},
})

// derpRegionErrorWarnable is a Warnable that warns the user that a DERP region is reporting an issue.
var derpRegionErrorWarnable = Register(&Warnable{
	Code:      "derp-region-error",
	Title:     "Relay server error",
	Severity:  SeverityMedium,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text: func(args Args) string {
		return fmt.Sprintf("The relay server #%v is reporting an issue: %v", args[ArgDERPRegionID], args[ArgError])
	},
})

// noUDP4BindWarnable is a Warnable that warns the user that Tailscale couldn't listen for incoming UDP connections.
var noUDP4BindWarnable = Register(&Warnable{
	Code:                "no-udp4-bind",
	Title:               "Incoming connections may fail",
	Severity:            SeverityHigh,
	DependsOn:           []*Warnable{NetworkStatusWarnable},
	Text:                StaticMessage("Tailscale couldn't listen for incoming UDP connections."),
	ImpactsConnectivity: true,
})

// mapResponseTimeoutWarnable is a Warnable that warns the user that Tailscale hasn't received a network map from the coordination server in a while.
var mapResponseTimeoutWarnable = Register(&Warnable{
	Code:      "mapresponse-timeout",
	Title:     "Network map response timeout",
	Severity:  SeverityMedium,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text: func(args Args) string {
		return fmt.Sprintf("Tailscale hasn't received a network map from the coordination server in %s.", args[ArgDuration])
	},
})

// tlsConnectionFailedWarnable is a Warnable that warns the user that Tailscale could not establish an encrypted connection with a server.
var tlsConnectionFailedWarnable = Register(&Warnable{
	Code:      "tls-connection-failed",
	Title:     "Encrypted connection failed",
	Severity:  SeverityMedium,
	DependsOn: []*Warnable{NetworkStatusWarnable},
	Text: func(args Args) string {
		return fmt.Sprintf("Tailscale could not establish an encrypted connection with '%q': %v", args[ArgServerName], args[ArgError])
	},
})

// magicsockReceiveFuncWarnable is a Warnable that warns the user that one of the Magicsock functions is not running.
var magicsockReceiveFuncWarnable = Register(&Warnable{
	Code:     "magicsock-receive-func-error",
	Title:    "MagicSock function not running",
	Severity: SeverityMedium,
	Text: func(args Args) string {
		return fmt.Sprintf("The MagicSock function %s is not running. You might experience connectivity issues.", args[ArgMagicsockFunctionName])
	},
})

// testWarnable is a Warnable that is used within this package for testing purposes only.
var testWarnable = Register(&Warnable{
	Code:     "test-warnable",
	Title:    "Test warnable",
	Severity: SeverityLow,
	Text: func(args Args) string {
		return args[ArgError]
	},
})

// applyDiskConfigWarnable is a Warnable that warns the user that there was an error applying the envknob config stored on disk.
var applyDiskConfigWarnable = Register(&Warnable{
	Code:     "apply-disk-config",
	Title:    "Could not apply configuration",
	Severity: SeverityMedium,
	Text: func(args Args) string {
		return fmt.Sprintf("An error occurred applying the Tailscale envknob configuration stored on disk: %v", args[ArgError])
	},
})

// controlHealthWarnable is a Warnable that warns the user that the coordination server is reporting an health issue.
var controlHealthWarnable = Register(&Warnable{
	Code:     "control-health",
	Title:    "Coordination server reports an issue",
	Severity: SeverityMedium,
	Text: func(args Args) string {
		return fmt.Sprintf("The coordination server is reporting an health issue: %v", args[ArgError])
	},
})
