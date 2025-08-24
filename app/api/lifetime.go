package api

import "context"

const (
	DevMode = "dev"

	ReleaseMode = "release"
)

// BeforeAppRunProcessor runs before app running
type BeforeAppRunProcessor interface {
	BeforeAppRun() error
}

// BeforeAppStopProcessor runs when app before stopping
type BeforeAppStopProcessor interface {
	BeforeAppStop(ctx context.Context) error
}

// AfterAppInitialProcessor calls AfterAppInit when app started(dependencies is resolved)
type AfterAppInitialProcessor interface {
	AfterAppInit() error
}
