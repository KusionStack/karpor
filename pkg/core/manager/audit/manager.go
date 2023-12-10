package audit

import (
	"context"
	"strings"

	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/scanner/kubeaudit"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

// AuditManager manages the auditing process of Kubernetes manifests using
// a KubeScanner.
type AuditManager struct {
	s scanner.KubeScanner // Interface for a Kubernetes scanner.
}

// NewAuditManager initializes a new instance of AuditManager with a KubeScanner.
func NewAuditManager() (*AuditManager, error) {
	// Create a new Kubernetes scanner instance.
	kubeauditScanner, err := kubeaudit.New()
	if err != nil {
		return nil, err
	}
	return &AuditManager{
		s: kubeauditScanner, // Set the scanner in the AuditManager.
	}, nil
}

// Audit performs a security audit on the provided manifest, returning a list
// of issues discovered during scanning.
func (m *AuditManager) Audit(ctx context.Context, manifest string) ([]*scanner.Issue, error) {
	// Retrieve logger from context and log the start of the audit.
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit of the specified manifest in AuditManager ...")

	// Execute the scan using the scanner's ScanManifest method.
	return m.s.ScanManifest(ctx, strings.NewReader(manifest))
}
