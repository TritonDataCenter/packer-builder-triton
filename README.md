## packer-builder-triton

[HashiCorp Packer](https://packer.io) builder for
[Joyent Triton](https://www.joyent.com).

To build from source:

1. Make sure [`glide`](https://github.com/Masterminds/glide) is installed (hint
   on OS X: `brew install glide`)
1. Set `GO15VENDOREXPERIMENT` to `1` if using Go 1.5 (default starting in 1.6)
1. `go get github.com/joyent/packer-builder-triton`
1. `cd $GOPATH/src/github.com/joyent/packer-builder-triton`
1. `make install`

### Example configuration

**NOTE:** Ensure that `SDC_KEY_ID` is in MD5 format (e.g.: `ssh-keygen -l -E md5
-f /Users/James/.ssh/joyent` - Mac OS X does not do this by default)

```json
{
    "variables": {
        "triton_url": "{{env `SDC_URL`}}",
        "triton_account": "{{env `SDC_ACCOUNT`}}",
        "triton_key_id": "{{env `SDC_KEY_ID`}}",
        "triton_key_path": "{{env `SDC_KEY_PATH`}}"
    },

    "builders": [
        {
            "type": "triton",
            "triton_url": "{{user `triton_url`}}",
            "triton_account": "{{user `triton_account`}}",
            "triton_key_id": "{{user `triton_key_id`}}",
            "triton_key_path": "{{user `triton_key_path`}}",

            "source_machine_package": "g3-standard-0.25-smartos",
            "source_machine_image": "842e6fa6-6e9b-11e5-8402-1b490459e334",
            "source_machine_networks": [
                "5983940e-58a5-4543-b732-c689b1fe4c08",
                "9ec60129-9034-47b4-b111-3026f9b1a10f"
            ],
            "source_machine_metadata": {
                "Key1": "Value1"
            },
            "source_machine_tags": {
                "Project": "Packer-Triton"
            },
            "source_machine_firewall_enabled": false,

            "image_name": "my-test-image",
            "image_version": "1.0.0",
            "image_description": "SDC Image created with Packer",
            "image_homepage": "http://jen20.com",
            "image_eula_url": "https://www.mozilla.org/media/MPL/2.0/index.815ca599c9df.txt",
            "image_acls": [],
            "image_tags": {
                "Project": "Packer-Triton"
            }
        }
    ],

    "provisioners": [
        {
            "type": "shell",
            "inline": ["touch foo.txt"]
        }
    ]
}
```
