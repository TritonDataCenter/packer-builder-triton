package triton

import (
	"fmt"
	"time"

	"github.com/joyent/gosdc/cloudapi"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// StepCreateImageFromMachine creates an image with the specified attributes
// from the machine with the given ID, and waits for the image to be created.
// The machine must be in the "stopped" state prior to this step being run.
type StepCreateImageFromMachine struct {
	ImageName        string
	ImageVersion     string
	ImageDescription string
	ImageHomepage    string
	ImageEULA        string
	ImageACL         []string
	ImageTags        map[string]string

	imageID string
}

func (s *StepCreateImageFromMachine) Run(state multistep.StateBag) multistep.StepAction {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	machineID := state.Get("machine").(string)

	ui.Say("Creating image from source machine...")
	opts := cloudapi.CreateImageFromMachineOpts{
		Machine:     machineID,
		Name:        s.ImageName,
		Version:     s.ImageVersion,
		Description: s.ImageDescription,
		Homepage:    s.ImageHomepage,
		EULA:        s.ImageEULA,
		ACL:         s.ImageACL,
		Tags:        s.ImageTags,
	}

	image, err := sdcClient.CreateImageFromMachine(opts)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem creating image from machine: %s", err))
		return multistep.ActionHalt
	}

	s.imageID = image.Id

	ui.Say("Waiting for image to become available...")
	err = waitForImageCreation(sdcClient, image.Id, 10*time.Minute)
	if err != nil {
		state.Put("error", fmt.Errorf("Problem waiting for image to become available: %s", err))
		return multistep.ActionHalt
	}

	state.Put("image", image.Id)

	return multistep.ActionContinue
}

func (s *StepCreateImageFromMachine) Cleanup(state multistep.StateBag) {
	sdcClient := state.Get("client").(*cloudapi.Client)
	ui := state.Get("ui").(packer.Ui)

	if s.imageID != "" {
		ui.Say("Deleting image...")
		err := sdcClient.DeleteImage(s.imageID)
		if err != nil {
			state.Put("error", fmt.Errorf("Problem deleting image: %s", err))
		}
	}
}
