~/.packer.d/plugins/packer-builder-triton: main.go triton/**
	@mkdir -p ~/.packer.d/plugins
	go build -o ~/.packer.d/plugins/packer-builder-triton .

vendor: glide.yaml glide.lock
	glide install

install: vendor ~/.packer.d/plugins/packer-builder-triton

test:
	go test $(glide nv)

.PHONY: install test
