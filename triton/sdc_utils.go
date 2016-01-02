package triton

import (
	"errors"
	"strings"
	"time"

	"github.com/joyent/gosdc/cloudapi"
)

// waitForMachineState uses the supplied client to wait for the state of
// the machine with the given ID to reach the state described in state.
// If timeout is reached before the machine reaches the required state, an
// error is returned. If the machine reaches the target state within the
// timeout, nil is returned.
func waitForMachineState(api *cloudapi.Client, id, state string, timeout time.Duration) error {
	return waitFor(
		func() (bool, error) {
			machine, err := api.GetMachine(id)
			return machine.State == state, err
		},
		3*time.Second,
		timeout,
	)
}

// waitForMachineDeletion uses the supplied client to wait for the machine
// with the given ID to be deleted. It is expected that the API call to delete
// the machine has already been issued at this point.
func waitForMachineDeletion(api *cloudapi.Client, id string, timeout time.Duration) error {
	return waitFor(
		func() (bool, error) {
			machine, err := api.GetMachine(id)
			if machine != nil {
				return false, nil
			}

			if err != nil {
				//TODO(jen20): is there a better way here than searching strings?
				if strings.Contains(err.Error(), "410") || strings.Contains(err.Error(), "404") {
					return true, nil
				}
			}

			return false, err
		},
		3*time.Second,
		timeout,
	)
}

func waitForImageCreation(api *cloudapi.Client, id string, timeout time.Duration) error {
	return waitFor(
		func() (bool, error) {
			image, err := api.GetImage(id)
			return image.OS != "", err
		},
		3*time.Second,
		timeout,
	)
}

func waitFor(f func() (bool, error), every, timeout time.Duration) error {
	start := time.Now()

	for time.Since(start) <= timeout {
		stop, err := f()
		if err != nil {
			return err
		}

		if stop {
			return nil
		}

		time.Sleep(every)
	}

	return errors.New("Timed out while waiting for resource change")
}
