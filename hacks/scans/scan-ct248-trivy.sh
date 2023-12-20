# https://wiki.ith.intel.com/display/SecTools/Trivy
# https://aquasecurity.github.io/trivy/v0.48/getting-started/installation/#debianubuntu-official
sudo trivy image localhost:5001/gardener-extension-cri-resmgr-installation-and-agent >/mnt/c/Users/ppalucki/Downloads/trivy-installation-image.txt
sudo trivy image localhost:5001/gardener-extension-cri-resmgr >/mnt/c/Users/ppalucki/Downloads/trivy-gecri-rm-image.txt
