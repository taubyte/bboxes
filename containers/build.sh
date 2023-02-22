#/bin/bash

cd containers

# Add an _builds directory if not already there
rm -r _builds
mkdir -p _builds

# Remove any fakeroot that may already exist
rm -r fakeroot
mkdir -p fakeroot/go
# mkdir -p fakeroot/{go,build} //{as,build}  || true

# Build go.tar
cp -r ./go/* ./common/* fakeroot/
tar cvf _builds/go.tar -C fakeroot/ .
rm -r fakeroot/*

# Build as.tar
cp -r ./as/* ./common/* fakeroot/
npm i --prefix fakeroot/utils/as-dir
tar cvf _builds/as.tar -C fakeroot/ .
rm -r fakeroot/*

# Build rs.tar 
cp -r ./rs/* ./common/* fakeroot/
tar cvf _builds/rs.tar -C fakeroot/ . 

# Clean fakeroot
rm -r fakeroot

exit $?


