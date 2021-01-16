package main

import (
	"crypto/sha512"
	"errors"
	"time"

	"github.com/google/uuid"
)

// HitService - Logic layer
type HitService struct {
	storage Storage
	salt    string
}

// NewHitService - Creates a new Hit Service
func NewHitService(storage Storage, salt string) HitService {
	return HitService{
		storage: storage,
		salt:    salt,
	}
}

// CountHit - Creates a new hit for a tenant's page
func (h *HitService) CountHit(tenantID string, pageName string, ip string) error {
	if tenantID == "" {
		return errors.New("Tenant id can not be blank")
	}

	if pageName == "" {
		return errors.New("Cannot use blank page name")
	}

	parsedTenantID, err := uuid.Parse(tenantID)

	if err != nil {
		return err
	}

	pageExists, err := h.storage.PageExists(parsedTenantID, pageName)

	if err != nil {
		return err
	}

	if pageExists {
		now := time.Now().UTC()
		bytes := []byte(ip + h.salt)
		hashed := sha512.New().Sum(bytes)
		footprint := string(hashed)

		err := h.storage.InsertHit(parsedTenantID, pageName, now, footprint)

		if err != nil {
			return err
		}
	} else {
		return errors.New("Page not found")
	}

	return nil
}
