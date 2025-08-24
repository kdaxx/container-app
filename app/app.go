package app

import (
	"context"
	"fmt"
	"github.com/kdaxx/container-app/app/api"
	api2 "github.com/kdaxx/container/v2/api"
	"github.com/kdaxx/container/v2/container"
	"go.uber.org/multierr"
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"
)

var c = container.NewContainer()

func Enable(modules []api2.BeanRegistrar) {
	c.ApplyRegistrars(modules)
}

func RunApplication() error {

	err := applyBeforeAppRunProcessors()
	if err != nil {
		return err
	}
	err = c.RunApplication()
	if err != nil {
		return err
	}

	err = applyInitializers()
	if err != nil {
		return err
	}

	return applyBeforeAppStopProcessors()

}

func applyBeforeAppRunProcessors() error {
	beforeRunners := c.GetBeanByType(reflect.TypeFor[api.BeforeAppRunProcessor]())
	for _, bean := range beforeRunners {
		beforeAppRun := bean.Bean().(api.BeforeAppRunProcessor)
		err := beforeAppRun.BeforeAppRun()
		if err != nil {
			return fmt.Errorf("%v run failed err:%v", reflect.TypeOf(beforeAppRun), err)
		}
	}
	return nil
}

func applyInitializers() error {
	appInitializers := c.GetBeanByType(reflect.TypeFor[api.AfterAppInitialProcessor]())
	for _, bean := range appInitializers {
		initializersBean := bean.Bean().(api.AfterAppInitialProcessor)
		err := initializersBean.AfterAppInit()
		if err != nil {
			return fmt.Errorf("%v init failed err:%v", reflect.TypeOf(initializersBean), err)
		}
	}
	return nil
}

func applyBeforeAppStopProcessors() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	var wait = 5
	log.Printf("app will be stopped in %d seconds\n", wait)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appStopRunners := c.GetBeanByType(reflect.TypeFor[api.BeforeAppStopProcessor]())

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	errs := make([]error, len(appStopRunners))

	// Concurrent execution appStopRunners to ensures that the application is closed within 5 seconds
	for _, bean := range appStopRunners {
		wg.Add(1)
		go func() {
			appStopRunner := bean.Bean().(api.BeforeAppStopProcessor)
			err := appStopRunner.BeforeAppStop(ctx)
			if err != nil {
				mutex.Lock()
				errs = append(errs, fmt.Errorf("%v run failed when app stopping:%v", reflect.TypeOf(appStopRunner), err))
				mutex.Unlock()
			}
			wg.Done()
		}()

	}
	wg.Wait()
	return multierr.Combine(errs...)
}
