CC=$ANDROID_HOME/ndk/27.0.12077973/toolchains/llvm/prebuilt/darwin-x86_64/bin/aarch64-linux-android21-clang CGO_ENABLED=1 GOOS=android GOARCH=arm64 go build main.go \
    && JAVA_SENDMSG/build.sh