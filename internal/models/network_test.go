package models

import (
	"testing"
)

func TestNewNetworkInfo(t *testing.T) {
	tests := []struct {
		name       string
		netName    string
		driver     string
		wantName   string
		wantDriver string
	}{
		{
			name:       "bridge network",
			netName:    "bridge",
			driver:     "bridge",
			wantName:   "bridge",
			wantDriver: "bridge",
		},
		{
			name:       "custom network with bridge driver",
			netName:    "frontend_net",
			driver:     "bridge",
			wantName:   "frontend_net",
			wantDriver: "bridge",
		},
		{
			name:       "overlay network",
			netName:    "swarm_network",
			driver:     "overlay",
			wantName:   "swarm_network",
			wantDriver: "overlay",
		},
		{
			name:       "host network",
			netName:    "host",
			driver:     "host",
			wantName:   "host",
			wantDriver: "host",
		},
		{
			name:       "macvlan network",
			netName:    "macvlan_net",
			driver:     "macvlan",
			wantName:   "macvlan_net",
			wantDriver: "macvlan",
		},
		{
			name:       "none network",
			netName:    "none",
			driver:     "none",
			wantName:   "none",
			wantDriver: "none",
		},
		{
			name:       "empty name",
			netName:    "",
			driver:     "bridge",
			wantName:   "",
			wantDriver: "bridge",
		},
		{
			name:       "empty driver",
			netName:    "test_net",
			driver:     "",
			wantName:   "test_net",
			wantDriver: "",
		},
		{
			name:       "both empty",
			netName:    "",
			driver:     "",
			wantName:   "",
			wantDriver: "",
		},
		{
			name:       "name with special characters",
			netName:    "my-network_v2.0",
			driver:     "bridge",
			wantName:   "my-network_v2.0",
			wantDriver: "bridge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNetworkInfo(tt.netName, tt.driver)

			if n == nil {
				t.Fatal("NewNetworkInfo returned nil")
			}

			if n.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", n.Name, tt.wantName)
			}

			if n.Driver != tt.wantDriver {
				t.Errorf("Driver = %q, want %q", n.Driver, tt.wantDriver)
			}
		})
	}
}

func TestNetworkInfo_DirectFieldAccess(t *testing.T) {
	// Test that the struct fields can be accessed and modified directly
	t.Run("create with direct initialization", func(t *testing.T) {
		n := &NetworkInfo{
			Name:   "custom_net",
			Driver: "bridge",
		}

		if n.Name != "custom_net" {
			t.Errorf("Name = %q, want %q", n.Name, "custom_net")
		}

		if n.Driver != "bridge" {
			t.Errorf("Driver = %q, want %q", n.Driver, "bridge")
		}
	})

	t.Run("modify fields directly", func(t *testing.T) {
		n := NewNetworkInfo("original", "bridge")

		n.Name = "modified"
		n.Driver = "overlay"

		if n.Name != "modified" {
			t.Errorf("Name = %q, want %q", n.Name, "modified")
		}

		if n.Driver != "overlay" {
			t.Errorf("Driver = %q, want %q", n.Driver, "overlay")
		}
	})
}

func TestNetworkInfo_ZeroValue(t *testing.T) {
	// Test zero-value struct behavior
	var n NetworkInfo

	if n.Name != "" {
		t.Errorf("zero-value Name = %q, want empty string", n.Name)
	}

	if n.Driver != "" {
		t.Errorf("zero-value Driver = %q, want empty string", n.Driver)
	}
}

func TestNetworkInfo_PointerVsValue(t *testing.T) {
	// NewNetworkInfo returns a pointer, verify behavior is consistent
	t.Run("pointer from constructor", func(t *testing.T) {
		n1 := NewNetworkInfo("test", "bridge")
		n2 := n1

		// Both should point to the same instance
		n2.Name = "modified"
		if n1.Name != "modified" {
			t.Error("pointer assignment should share the same instance")
		}
	})

	t.Run("value copy behavior", func(t *testing.T) {
		n1 := NewNetworkInfo("test", "bridge")
		n2 := *n1 // Dereference to create a value copy

		n2.Name = "modified"
		if n1.Name == "modified" {
			t.Error("value copy should not affect original")
		}

		if n1.Name != "test" {
			t.Errorf("original Name = %q, want %q", n1.Name, "test")
		}

		if n2.Name != "modified" {
			t.Errorf("copy Name = %q, want %q", n2.Name, "modified")
		}
	})
}
