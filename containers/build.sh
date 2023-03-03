#/bin/bash

cd containers

# TODO: Further Optimize this code  

# Add an _builds directory if not already there
rm -r _builds
mkdir -p _builds

# Remove any fakeroot that may already exist
rm -r fakeroot
mkdir -p fakeroot

# # Build go.tar
cp -r ./go/* ./common/* fakeroot/
tar cvf _builds/go.tar -C fakeroot/ .
rm -r fakeroot/*

# Build go_test_examples.tar 
cp -r ./test-examples/go/* ./test-examples/common/* fakeroot/
tar cvf _builds/go_test_examples.tar -C fakeroot/ . 
rm -r fakeroot/*

# # Build as.tar
cp -r ./as/* ./common/* fakeroot/
npm i --prefix fakeroot/utils/as-dir
tar cvf _builds/as.tar -C fakeroot/ .
rm -r fakeroot/*

# # Build as_test_examples.tar
cp -r ./test-examples/as/* ./test-examples/common/* fakeroot/
npm i --prefix fakeroot/utils/as-dir
tar cvf _builds/as_test_examples.tar -C fakeroot/ .
rm -r fakeroot/*

# # Build rs.tar 
cp -r ./rs/* ./common/* fakeroot/
tar cvf _builds/rs.tar -C fakeroot/ . 
rm -r fakeroot/*

# # Build rs_test_examples.tar 
cp -r ./test-examples/rs/* ./test-examples/common/* fakeroot/
tar cvf _builds/rs_test_examples.tar -C fakeroot/ . 

# # Clean fakeroot
rm -r fakeroot

exit $?


