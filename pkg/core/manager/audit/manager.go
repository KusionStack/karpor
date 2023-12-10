package audit

import (
	"context"
	"strings"

	"github.com/KusionStack/karbour/pkg/scanner"
	"github.com/KusionStack/karbour/pkg/scanner/kubeaudit"
	"github.com/KusionStack/karbour/pkg/util/ctxutil"
)

type AuditManager struct {
	s scanner.KubeScanner
}

func NewAuditManager() (*AuditManager, error) {
	kubeauditScanner, err := kubeaudit.New()
	if err != nil {
		return nil, err
	}

	return &AuditManager{
		s: kubeauditScanner,
	}, nil
}

func (m *AuditManager) Audit(ctx context.Context, manifest string) ([]*scanner.Issue, error) {
	log := ctxutil.GetLogger(ctx)
	log.Info("Starting audit the specified manifest in AuditManager ...")
	return m.s.ScanManifest(ctx, strings.NewReader(manifest))
}
