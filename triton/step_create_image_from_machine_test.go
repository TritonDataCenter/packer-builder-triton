package triton

import (
	"errors"
	"testing"

	"github.com/mitchellh/multistep"
)

func TestStepCreateImageFromMachine(t *testing.T) {
	state := testState(t)
	step := new(StepCreateImageFromMachine)
	defer step.Cleanup(state)

	driver := state.Get("driver").(*DriverMock)
	state.Put("machine", "test-machine-id")

	if action := step.Run(state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}

	imageIdRaw, ok := state.GetOk("image")
	if !ok {
		t.Fatalf("should have image")
	}

	step.Cleanup(state)

	if driver.DeleteImageId != imageIdRaw.(string) {
		t.Fatalf("should've deleted image (%s != %s)", driver.DeleteImageId, imageIdRaw.(string))
	}
}

func TestStepCreateImageFromMachine_CreateImageFromMachineError(t *testing.T) {
	state := testState(t)
	step := new(StepCreateImageFromMachine)
	defer step.Cleanup(state)

	driver := state.Get("driver").(*DriverMock)
	state.Put("machine", "test-machine-id")

	driver.CreateImageFromMachineErr = errors.New("error")

	if action := step.Run(state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %#v", action)
	}

	if _, ok := state.GetOk("error"); !ok {
		t.Fatalf("should have error")
	}

	if _, ok := state.GetOk("image"); ok {
		t.Fatalf("should NOT have image")
	}
}

func TestStepCreateImageFromMachine_WaitForImageCreationError(t *testing.T) {
	state := testState(t)
	step := new(StepCreateImageFromMachine)
	defer step.Cleanup(state)

	driver := state.Get("driver").(*DriverMock)
	state.Put("machine", "test-machine-id")

	driver.WaitForImageCreationErr = errors.New("error")

	if action := step.Run(state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %#v", action)
	}

	if _, ok := state.GetOk("error"); !ok {
		t.Fatalf("should have error")
	}

	if _, ok := state.GetOk("image"); ok {
		t.Fatalf("should NOT have image")
	}
}

func TestStepCreateImageFromMachine_DeleteImageError(t *testing.T) {
	state := testState(t)
	step := new(StepCreateImageFromMachine)
	defer step.Cleanup(state)

	driver := state.Get("driver").(*DriverMock)
	state.Put("machine", "test-machine-id")

	if action := step.Run(state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}

	_, ok := state.GetOk("image")
	if !ok {
		t.Fatalf("should have image")
	}

	driver.DeleteImageErr = errors.New("error")
	step.Cleanup(state)

	if _, ok := state.GetOk("error"); !ok {
		t.Fatalf("should have error")
	}

	if _, ok := state.GetOk("image"); !ok {
		t.Fatalf("should have image")
	}
}
