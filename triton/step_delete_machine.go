package triton

import (
	"fmt"
	"time"

	"github.com/joyent/gosdc/cloudapi"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepDeleteMachine deletes the machine with the ID specified in state["machine"]
type StepDeleteMachine struct{}

func (s *StepDeleteMachine) Run(state multistep.StateBag) multistep.StepAction {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	machineID := state.Get("machine").(string)

	ui.Say("Deleting source machine...")
	err := sdcClient.DeleteMachine(machineID)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem deleting source machine: %s", err))
		return multistep.ActionHalt
	}

	ui.Say("Waiting for source machine to be deleted...")
	err = waitForMachineDeletion(sdcClient, machineID, 10*time.Minute)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem waiting for source machine to be deleted: %s", err))
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepDeleteMachine) Cleanup(state multistep.StateBag) {
	// No clean up to do here...
}
