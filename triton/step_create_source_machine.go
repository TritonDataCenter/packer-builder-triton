package triton

import (
	"fmt"
	"time"

	"github.com/joyent/gosdc/cloudapi"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
)

// StepCreateSourceMachine creates an machine with the specified attributes
// and waits for it to become available for provisioners.
type StepCreateSourceMachine struct {
	MachineName            string
	MachinePackage         string
	MachineImage           string
	MachineNetworks        []string
	MachineMetadata        map[string]string
	MachineTags            map[string]string
	MachineFirewallEnabled bool

	Comm communicator.Config `mapstructure:",squash"`

	machineId string
}

func (s *StepCreateSourceMachine) Run(state multistep.StateBag) multistep.StepAction {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating source machine...")
	opts := cloudapi.CreateMachineOpts{
		Package:         s.MachinePackage,
		Image:           s.MachineImage,
		Networks:        s.MachineNetworks,
		Metadata:        s.MachineMetadata,
		Tags:            s.MachineTags,
		FirewallEnabled: s.MachineFirewallEnabled,
	}

	// Only supply a name if it was supplied in the template; the SDC API will
	// assign a default random name otherwise.
	if s.MachineName != "" {
		opts.Name = s.MachineName
	}

	machine, err := sdcClient.CreateMachine(opts)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem creating source machine: %s", err))
		return multistep.ActionHalt
	}

	s.machineId = machine.Id

	ui.Say("Waiting for source machine to become available...")
	err = waitForMachineState(sdcClient, machine.Id, "running", 10*time.Minute)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem waiting for source machine to become available: %s", err))
		return multistep.ActionHalt
	}

	state.Put("machine", machine.Id)

	return multistep.ActionContinue
}

func (s *StepCreateSourceMachine) Cleanup(state multistep.StateBag) {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	if s.machineId != "" {
		ui.Say(fmt.Sprintf("Stopping source machine (%s)...", s.machineId))
		err := sdcClient.StopMachine(s.machineId)
		if err != nil {
			state.Put("error", fmt.Errorf("Problem stopping source machine: %s", err))
			return
		}

		ui.Say(fmt.Sprintf("Waiting for source machine to stop (%s)...", s.machineId))
		err = waitForMachineState(sdcClient, s.machineId, "stopped", 10*time.Minute)
		if err != nil {
			state.Put("error", fmt.Errorf("Problem waiting for source machine to stop: %s", err))
			return
		}

		ui.Say("Deleting source machine...")
		err = sdcClient.DeleteMachine(s.machineId)
		if err != nil {
			state.Put("error", fmt.Errorf("Problem deleting source machine: %s", err))
			return
		}
	}
}
