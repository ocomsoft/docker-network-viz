package models

import (
	"testing"
)

func TestNewContainerInfo(t *testing.T) {
	tests := []struct {
		name          string
		containerName string
	}{
		{
			name:          "simple name",
			containerName: "web_app",
		},
		{
			name:          "name with dashes",
			containerName: "my-web-app",
		},
		{
			name:          "empty name",
			containerName: "",
		},
		{
			name:          "name with numbers",
			containerName: "app123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContainerInfo(tt.containerName)

			if c == nil {
				t.Fatal("NewContainerInfo returned nil")
			}

			if c.Name != tt.containerName {
				t.Errorf("Name = %q, want %q", c.Name, tt.containerName)
			}

			if c.Aliases == nil {
				t.Error("Aliases should not be nil")
			}

			if len(c.Aliases) != 0 {
				t.Errorf("Aliases length = %d, want 0", len(c.Aliases))
			}

			if c.Networks == nil {
				t.Error("Networks should not be nil")
			}

			if len(c.Networks) != 0 {
				t.Errorf("Networks length = %d, want 0", len(c.Networks))
			}
		})
	}
}

func TestContainerInfo_AddAlias(t *testing.T) {
	t.Run("add new alias", func(t *testing.T) {
		c := NewContainerInfo("test")

		added := c.AddAlias("web")
		if !added {
			t.Error("AddAlias should return true for new alias")
		}

		if len(c.Aliases) != 1 {
			t.Errorf("Aliases length = %d, want 1", len(c.Aliases))
		}

		if c.Aliases[0] != "web" {
			t.Errorf("Alias = %q, want %q", c.Aliases[0], "web")
		}
	})

	t.Run("add duplicate alias", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddAlias("web")

		added := c.AddAlias("web")
		if added {
			t.Error("AddAlias should return false for duplicate alias")
		}

		if len(c.Aliases) != 1 {
			t.Errorf("Aliases length = %d, want 1", len(c.Aliases))
		}
	})

	t.Run("add multiple aliases", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddAlias("web")
		c.AddAlias("api")
		c.AddAlias("app")

		if len(c.Aliases) != 3 {
			t.Errorf("Aliases length = %d, want 3", len(c.Aliases))
		}
	})
}

func TestContainerInfo_AddNetwork(t *testing.T) {
	t.Run("add new network", func(t *testing.T) {
		c := NewContainerInfo("test")

		added := c.AddNetwork("bridge")
		if !added {
			t.Error("AddNetwork should return true for new network")
		}

		if len(c.Networks) != 1 {
			t.Errorf("Networks length = %d, want 1", len(c.Networks))
		}

		if c.Networks[0] != "bridge" {
			t.Errorf("Network = %q, want %q", c.Networks[0], "bridge")
		}
	})

	t.Run("add duplicate network", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddNetwork("bridge")

		added := c.AddNetwork("bridge")
		if added {
			t.Error("AddNetwork should return false for duplicate network")
		}

		if len(c.Networks) != 1 {
			t.Errorf("Networks length = %d, want 1", len(c.Networks))
		}
	})

	t.Run("add multiple networks", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddNetwork("bridge")
		c.AddNetwork("frontend")
		c.AddNetwork("backend")

		if len(c.Networks) != 3 {
			t.Errorf("Networks length = %d, want 3", len(c.Networks))
		}
	})
}

func TestContainerInfo_HasNetwork(t *testing.T) {
	c := NewContainerInfo("test")
	c.AddNetwork("bridge")
	c.AddNetwork("frontend")

	tests := []struct {
		name    string
		network string
		want    bool
	}{
		{
			name:    "existing network",
			network: "bridge",
			want:    true,
		},
		{
			name:    "another existing network",
			network: "frontend",
			want:    true,
		},
		{
			name:    "non-existing network",
			network: "backend",
			want:    false,
		},
		{
			name:    "empty string",
			network: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.HasNetwork(tt.network)
			if got != tt.want {
				t.Errorf("HasNetwork(%q) = %v, want %v", tt.network, got, tt.want)
			}
		})
	}
}

func TestContainerInfo_HasAlias(t *testing.T) {
	c := NewContainerInfo("test")
	c.AddAlias("web")
	c.AddAlias("api")

	tests := []struct {
		name  string
		alias string
		want  bool
	}{
		{
			name:  "existing alias",
			alias: "web",
			want:  true,
		},
		{
			name:  "another existing alias",
			alias: "api",
			want:  true,
		},
		{
			name:  "non-existing alias",
			alias: "app",
			want:  false,
		},
		{
			name:  "empty string",
			alias: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.HasAlias(tt.alias)
			if got != tt.want {
				t.Errorf("HasAlias(%q) = %v, want %v", tt.alias, got, tt.want)
			}
		})
	}
}

func TestContainerInfo_SortedNetworks(t *testing.T) {
	t.Run("returns sorted copy", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddNetwork("zebra")
		c.AddNetwork("alpha")
		c.AddNetwork("beta")

		sorted := c.SortedNetworks()

		// Check sorted order
		expected := []string{"alpha", "beta", "zebra"}
		for i, want := range expected {
			if sorted[i] != want {
				t.Errorf("sorted[%d] = %q, want %q", i, sorted[i], want)
			}
		}

		// Verify original is unchanged
		if c.Networks[0] != "zebra" {
			t.Error("SortedNetworks should not modify original slice")
		}
	})

	t.Run("empty networks", func(t *testing.T) {
		c := NewContainerInfo("test")
		sorted := c.SortedNetworks()

		if len(sorted) != 0 {
			t.Errorf("SortedNetworks length = %d, want 0", len(sorted))
		}
	})
}

func TestContainerInfo_SortedAliases(t *testing.T) {
	t.Run("returns sorted copy", func(t *testing.T) {
		c := NewContainerInfo("test")
		c.AddAlias("zebra")
		c.AddAlias("alpha")
		c.AddAlias("beta")

		sorted := c.SortedAliases()

		// Check sorted order
		expected := []string{"alpha", "beta", "zebra"}
		for i, want := range expected {
			if sorted[i] != want {
				t.Errorf("sorted[%d] = %q, want %q", i, sorted[i], want)
			}
		}

		// Verify original is unchanged
		if c.Aliases[0] != "zebra" {
			t.Error("SortedAliases should not modify original slice")
		}
	})

	t.Run("empty aliases", func(t *testing.T) {
		c := NewContainerInfo("test")
		sorted := c.SortedAliases()

		if len(sorted) != 0 {
			t.Errorf("SortedAliases length = %d, want 0", len(sorted))
		}
	})
}

func TestContainerInfo_NetworkCount(t *testing.T) {
	tests := []struct {
		name     string
		networks []string
		want     int
	}{
		{
			name:     "no networks",
			networks: []string{},
			want:     0,
		},
		{
			name:     "one network",
			networks: []string{"bridge"},
			want:     1,
		},
		{
			name:     "multiple networks",
			networks: []string{"bridge", "frontend", "backend"},
			want:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContainerInfo("test")
			for _, n := range tt.networks {
				c.AddNetwork(n)
			}

			got := c.NetworkCount()
			if got != tt.want {
				t.Errorf("NetworkCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestContainerInfo_AliasCount(t *testing.T) {
	tests := []struct {
		name    string
		aliases []string
		want    int
	}{
		{
			name:    "no aliases",
			aliases: []string{},
			want:    0,
		},
		{
			name:    "one alias",
			aliases: []string{"web"},
			want:    1,
		},
		{
			name:    "multiple aliases",
			aliases: []string{"web", "api", "app"},
			want:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewContainerInfo("test")
			for _, a := range tt.aliases {
				c.AddAlias(a)
			}

			got := c.AliasCount()
			if got != tt.want {
				t.Errorf("AliasCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestContainerInfo_Clone(t *testing.T) {
	t.Run("creates deep copy", func(t *testing.T) {
		original := NewContainerInfo("original")
		original.AddAlias("web")
		original.AddAlias("api")
		original.AddNetwork("bridge")
		original.AddNetwork("frontend")

		clone := original.Clone()

		// Verify values are equal
		if clone.Name != original.Name {
			t.Errorf("Clone Name = %q, want %q", clone.Name, original.Name)
		}

		if len(clone.Aliases) != len(original.Aliases) {
			t.Errorf("Clone Aliases length = %d, want %d", len(clone.Aliases), len(original.Aliases))
		}

		if len(clone.Networks) != len(original.Networks) {
			t.Errorf("Clone Networks length = %d, want %d", len(clone.Networks), len(original.Networks))
		}

		// Verify it's a different instance
		if clone == original {
			t.Error("Clone should return a different pointer")
		}
	})

	t.Run("modifications don't affect original", func(t *testing.T) {
		original := NewContainerInfo("original")
		original.AddAlias("web")
		original.AddNetwork("bridge")

		clone := original.Clone()

		// Modify clone
		clone.Name = "modified"
		clone.AddAlias("new-alias")
		clone.AddNetwork("new-network")

		// Verify original is unchanged
		if original.Name != "original" {
			t.Errorf("Original Name changed to %q", original.Name)
		}

		if len(original.Aliases) != 1 {
			t.Errorf("Original Aliases length changed to %d", len(original.Aliases))
		}

		if len(original.Networks) != 1 {
			t.Errorf("Original Networks length changed to %d", len(original.Networks))
		}
	})

	t.Run("clone empty container", func(t *testing.T) {
		original := NewContainerInfo("empty")
		clone := original.Clone()

		if clone.Name != "empty" {
			t.Errorf("Clone Name = %q, want %q", clone.Name, "empty")
		}

		if len(clone.Aliases) != 0 {
			t.Errorf("Clone Aliases length = %d, want 0", len(clone.Aliases))
		}

		if len(clone.Networks) != 0 {
			t.Errorf("Clone Networks length = %d, want 0", len(clone.Networks))
		}
	})
}

func TestContainerInfo_DirectFieldAccess(t *testing.T) {
	// Test that the struct fields can be accessed directly
	c := &ContainerInfo{
		Name:     "direct",
		Aliases:  []string{"a1", "a2"},
		Networks: []string{"n1", "n2"},
	}

	if c.Name != "direct" {
		t.Errorf("Name = %q, want %q", c.Name, "direct")
	}

	if len(c.Aliases) != 2 {
		t.Errorf("Aliases length = %d, want 2", len(c.Aliases))
	}

	if len(c.Networks) != 2 {
		t.Errorf("Networks length = %d, want 2", len(c.Networks))
	}
}
