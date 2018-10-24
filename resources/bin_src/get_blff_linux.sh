#!/usr/bin/env bash
export GOPATH=$(go env GOPATH)

# Downloading compiled blff executables on linux just to get a successful qclauncher build with the packaged bin resources -- building requires Windows (Visual Studio)
# If you still want the blff source: git clone https://github.com/syncore/blff.git

cd $GOPATH/src/github.com/syncore/qclauncher
curl https://www.syncore.org/blff-v0.1.zip -o $GOPATH/src/github.com/syncore/qclauncher/blff-v0.1.zip
unzip -o $GOPATH/src/github.com/syncore/qclauncher/blff-v0.1.zip
rm -rf $GOPATH/src/github.com/syncore/qclauncher/blff-v0.1.zip