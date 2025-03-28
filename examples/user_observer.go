package examples

import (
	"fmt"
	"godm/pkg/odm"
)

type UserObserver struct {
	odm.ModelObserver
}

func (UserObserver) Creating(model interface{}) error {
	fmt.Println("[Observer] Creating:", model)
	return nil
}

func (UserObserver) Created(model interface{}) error {
	fmt.Println("[Observer] Created:", model)
	return nil
}

func (UserObserver) Deleted(model interface{}) error {
	fmt.Println("[Observer] Deleted:", model)
	return nil
}

func (UserObserver) Updating(model interface{}) error {
	return nil
}

func (UserObserver) Updated(model interface{}) error {
	return nil
}

func (UserObserver) InterestedIn(stage string) bool {
	return stage == "creating" || stage == "created" || stage == "deleted" || stage == "updating" || stage == "updated"
}

func (UserObserver) Accepts(model interface{}) bool {
	_, ok := model.(*User)
	return ok
}

func (UserObserver) Priority() int {
	return 100
}
