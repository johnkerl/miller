# Failed attempts to create a snap interactively

2026-01-02 I used an Ubuntu 24.04 EC2 instance. I followed https://documentation.ubuntu.com/snapcraft/stable/. Error messages said things like

```
A network related operation failed in a context of no network access.
Recommended resolution: Verify that the environment has internet connectivity; see https://canonical-craft-providers.readthedocs-hosted.com/en/latest/explanation/ for further reference.
Full execution log: '/home/ubuntu/.local/state/snapcraft/log/snapcraft-20260102-170252.488632.log'
```

when there was in fact no network problem. I remained confused.

```
$ sudo snapcraft pack

$ lxc list

$ snapcraft pack --destructive-mode

$ snapcraft pack --use-multipass

$ sudo snap install multipass

$ snapcraft pack --use-multipass

$ sudo lxd init --auto

$ lxc network list

$ sudo snapcraft pack

$ sudo snap set snapcraft provider=multipass

$ sudo snapcraft pack --destructive-mode

[This created miller_6.15.0_arm64.snap]

$ snapcraft upload --release=stable *.snap
No keyring found to store or retrieve credentials from.
Recommended resolution: Ensure the keyring is working or SNAPCRAFT_STORE_CREDENTIALS is correctly exported into the environment
For more information, check out: https://documentation.ubuntu.com/snapcraft/stable/how-to/publishing/authenticate
Full execution log: '/home/ubuntu/.local/state/snapcraft/log/snapcraft-20260102-172357.599171.log'

$ ll *.snap
-rw-r--r-- 1 root root 8994816 Jan  2 17:22 miller_6.15.0_arm64.snap

$ snap install *.snap
error: access denied (try with sudo)

$ sudo snap install *.snap
error: cannot find signatures with metadata for snap/component "miller_6.15.0_arm64.snap"
```

Conclusion:

* I got cryptic error messages with various permutations.
* Through trial and error I got a `.snap` file with `sudo` and `multipass` and `--destructive-mode`.
* Even then, I got a `.snap` file only for the current machine's arch, and the resulting `.snap` file was not locally installable.
* This led me to try a GitHub Action.

# Info from Claude about auto-releasing

Here's how to set up automatic Snap publishing from GitHub releases:

## 1. Create snapcraft.yaml

First, ensure you have a proper `snapcraft.yaml` in your repo root (or in a `snap/` directory):

```yaml
name: your-app-name
base: core22  # or core24 for Ubuntu 24.04
version: git  # automatically uses git tags
summary: Single-line summary
description: |
  Longer description of your application

grade: stable  # or devel
confinement: strict  # or classic, devmode

apps:
  your-app-name:
    command: bin/your-binary
    plugs:
      - home
      - network

parts:
  your-app:
    plugin: nil  # change based on your build system (go, python, etc.)
    source: .
    # Add build steps as needed
```

## 2. Get Snapcraft credentials

Export your Snapcraft login credentials:

```bash
snapcraft export-login --snaps=miller --channels=stable,candidate,beta,edge snapcraft-token.txt
```

This creates a token file with limited permissions for just your snap.

## 3. Add token to GitHub Secrets

1. Go to your GitHub repo → Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `SNAPCRAFT_TOKEN`
4. Value: Paste the entire contents of `snapcraft-token.txt`

## 4. Create GitHub Action workflow

Create `.github/workflows/release.yml`:

```yaml
name: Release to Snap Store

on:
  release:
    types: [published]

jobs:
  snap:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Build snap
        uses: snapcore/action-build@v1
        id: build

      - name: Publish to Snap Store
        uses: snapcore/action-publish@v1
        env:
          SNAPCRAFT_STORE_CREDENTIALS: ${{ secrets.SNAPCRAFT_TOKEN }}
        with:
          snap: ${{ steps.build.outputs.snap }}
          # release: stable  # or edge, beta, candidate
          release: edge
```

## Tips

- **Version handling**: Using `version: git` in snapcraft.yaml automatically uses your git tag as the version
- **Channels**: Start with `edge` channel for testing, then promote to `stable` once confident
- **Multiple architectures**: Add a build matrix if you need to support arm64, etc.
- **Testing before stable**: Consider publishing to `candidate` or `beta` first, then manually promote to `stable` after testing

Now when you create a GitHub release with a tag (e.g., `v1.0.0`), the workflow will automatically build and publish your snap!
