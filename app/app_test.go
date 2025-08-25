package app

import (
	"github.com/kdaxx/container/v2/api"
	"testing"
	"time"
)

func TestRunApplication(t *testing.T) {
	Enable([]api.BeanRegistrar{
		NewRegistrar(),
	})

	go func() {
		err := RunApplication()
		if err != nil {
			t.Errorf("application failed to run: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)
	StopApplication()
}

func TestVersion(t *testing.T) {
	t.Log(api.VERSION)
}
