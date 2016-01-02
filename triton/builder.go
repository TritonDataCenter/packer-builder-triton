package triton

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
)

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	ui.Say(fmt.Sprintf("%# v", pretty.Formatter(b.config)))
	ui.Say(b.config.Endpoint)
	return nil, nil
}

// Cancel cancels a possibly running Builder. This should block until
// the builder actually cancels and cleans up after itself.
func (b *Builder) Cancel() {

}
