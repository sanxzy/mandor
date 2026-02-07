package domain

import "time"

const (
	BriefStatusDraft    = "draft"
	BriefStatusActive   = "active"
	BriefStatusArchived = "archived"
)

// Capability represents a new or modified capability in a Brief
type Capability struct {
	ID          string `json:"id"`          // capability-id (derived from name via ToSlug)
	Name        string `json:"name"`        // Display name
	Description string `json:"description"` // Description of what capability provides
}

// Brief represents the root intent document for a backlog
type Brief struct {
	ID                   string       `json:"id"`                    // brief-id (derived from name via ToSlug)
	BacklogID            string       `json:"backlog_id"`            // Reference to backlog
	Status               string       `json:"status"`                // draft | active | archived
	Why                  string       `json:"why"`                   // Problem statement, 100-5000 chars
	WhatChanges          []string     `json:"what_changes"`          // List of capability descriptions
	Impact               BriefImpact  `json:"impact"`                // Technical stack, affected systems, dependencies
	NewCapabilities      []Capability `json:"new_capabilities"`      // New capabilities being added
	ModifiedCapabilities []Capability `json:"modified_capabilities"` // Existing capabilities being modified
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	CreatedBy            string       `json:"created_by"`
	UpdatedBy            string       `json:"updated_by"`
}

// BriefImpact describes the technical and systemic impact
type BriefImpact struct {
	TechnicalStack  []string `json:"technical_stack"`  // Technology choices
	AffectedSystems []string `json:"affected_systems"` // Systems impacted
	Dependencies    []string `json:"dependencies"`     // External dependencies
}

// BriefCreateInput represents input for creating a Brief
type BriefCreateInput struct {
	BacklogID    string
	Name         string            // Brief name (converted to ID via ToSlug)
	Why          string            // 100-5000 chars
	WhatChanges  []string          // Capability descriptions
	Impact       *BriefImpact      // Optional, defaults to empty
	Capabilities []CapabilityInput // New and modified capabilities
}

// CapabilityInput represents input for a capability
type CapabilityInput struct {
	Name        string // Capability name (converted to ID via ToSlug)
	Description string
	Modified    bool // true if modifying existing, false if new
}

// BriefUpdateInput represents input for updating a Brief
type BriefUpdateInput struct {
	ID                   string
	Status               *string
	Why                  *string
	WhatChanges          *[]string
	Impact               *BriefImpact
	NewCapabilities      *[]Capability
	ModifiedCapabilities *[]Capability
}

// ValidateBriefWhy validates why section length (100-5000 chars)
func ValidateBriefWhy(why string) bool {
	return len(why) >= 100 && len(why) <= 5000
}

// ValidateCapabilityID validates capability ID format [a-z0-9-]+
func ValidateCapabilityID(id string) bool {
	if len(id) == 0 {
		return false
	}
	for _, c := range id {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}
