package gen

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID(version int, namespace, name string) (string, error) {
	switch version {
	case 1:
		id, err := uuid.NewUUID()
		if err != nil {
			return "", fmt.Errorf("failed to generate v1 UUID: %w", err)
		}
		return id.String(), nil
	case 4:
		id, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("failed to generate v4 UUID: %w", err)
		}
		return id.String(), nil
	case 5:
		ns, err := resolveNamespace(namespace)
		if err != nil {
			return "", err
		}
		id := uuid.NewSHA1(ns, []byte(name))
		return id.String(), nil
	case 7:
		id, err := uuid.NewV7()
		if err != nil {
			return "", fmt.Errorf("failed to generate v7 UUID: %w", err)
		}
		return id.String(), nil
	default:
		return "", fmt.Errorf("unsupported UUID version: %d (supported: 1, 4, 5, 7)", version)
	}
}

func resolveNamespace(ns string) (uuid.UUID, error) {
	switch strings.ToLower(ns) {
	case "dns":
		return uuid.NameSpaceDNS, nil
	case "url":
		return uuid.NameSpaceURL, nil
	case "oid":
		return uuid.NameSpaceOID, nil
	case "x500":
		return uuid.NameSpaceX500, nil
	default:
		// Try parsing as UUID.
		id, err := uuid.Parse(ns)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("invalid namespace %q: must be dns, url, oid, x500, or a valid UUID", ns)
		}
		return id, nil
	}
}
