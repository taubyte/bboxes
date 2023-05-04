#/bin/bash

cd containers

# TODO: This is really ugly, should optimize

# Add an _builds directory if not already there
rm -r _builds
mkdir -p _builds/production
mkdir -p _builds/test_examples

# Remove any fakeroot that may already exist
rm -r fakeroot
mkdir -p fakeroot

# # Build go.tar
cp -r ./go/* ./common/* fakeroot/
tar cvf _builds/production/go.tar -C fakeroot/ .
rm -r fakeroot/*

# Build go_test_examples.tar 
cp -r ./test-examples/go/* ./test-examples/common/* fakeroot/
tar cvf _builds/test_examples/go.tar -C fakeroot/ . 
rm -r fakeroot/*

# # Build as.tar
cp -r ./as/* ./common/* fakeroot/
npm i --prefix fakeroot/utils/as-dir
tar cvf _builds/production/as.tar -C fakeroot/ .
rm -r fakeroot/*

# # Build as_test_examples.tar
cp -r ./test-examples/as/* ./test-examples/common/* fakeroot/
npm i --prefix fakeroot/utils/as-dir
tar cvf _builds/test_examples/as.tar -C fakeroot/ .
rm -r fakeroot/*

# # Build rs.tar 
cp -r ./rs/* ./common/* fakeroot/
tar cvf _builds/production/rs.tar -C fakeroot/ . 
rm -r fakeroot/*

# # Build rs_test_examples.tar 
cp -r ./test-examples/rs/* ./test-examples/common/* fakeroot/
tar cvf _builds/test_examples/rs.tar -C fakeroot/ . 

# # Clean fakeroot
rm -r fakeroot

exit $?


