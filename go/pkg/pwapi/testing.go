package pwapi

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"pathwar.land/go/pkg/pwdb"
	"pathwar.land/go/pkg/pwsso"
)

func TestingService(t *testing.T, opts ServiceOpts) (Service, func()) {
	t.Helper()

	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	db := pwdb.TestingSqliteDB(t, opts.Logger)
	sso := pwsso.TestingSSO(t, opts.Logger)

	api, err := NewService(db, sso, opts)
	if err != nil {
		t.Fatalf("init api: %v", err)
	}

	cleanup := func() {
		api.Close()
		db.Close()
	}

	return api, cleanup
}

func TestingServer(t *testing.T, ctx context.Context, opts ServerOpts) (*Server, func()) {
	t.Helper()

	svc, svcCleanup := TestingService(t, ServiceOpts{Logger: opts.Logger})

	if opts.HTTPBind == "" {
		opts.HTTPBind = "127.0.0.1:0"
	}
	if opts.GRPCBind == "" {
		opts.GRPCBind = "127.0.0.1:0"
	}

	server, err := NewServer(ctx, svc, opts)
	if err != nil {
		t.Fatalf("init server: %v", err)
	}

	cleanup := func() {
		server.Close()
		svcCleanup()
	}

	go func() {
		if err := server.Run(); err != nil {
			t.Logf("server shutdown, err: %v", err)
		}
	}()

	return server, cleanup
}