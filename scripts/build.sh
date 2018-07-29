#!/usr/bin/env bash

platforms=("linux/amd64" "linux/386" "windows/amd64" "windows/386" "darwin/amd64")

# Cd to root of the project
pathToScript="`dirname \"$0\"`"
cd ${pathToScript}/..

# Run tests
go test ./...
if [[ $? != 0 ]]; then
    (>&2 echo "Running tests failed.")
    exit $?
fi

for platform in "${platforms[@]}"
do
    echo "Building for ${platform}..."

    platformSplit=(${platform//\// })
    GOOS=${platformSplit[0]}
    GOARCH=${platformSplit[1]}
    buildDir="build/"${GOOS}'-'${GOARCH}

    # Recreate buildDir
    rm -r ${buildDir} 2> /dev/null
    mkdir -p ${buildDir}

    # Copy assets dir and rename it to data
    cp -r assets ${buildDir}
    mv ${buildDir}/assets  ${buildDir}/data

    execName="lia"
    if [ ${GOOS} = "windows" ]; then
        execName+='.exe'
    fi

    env GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${buildDir}/${execName} cmd/lia/main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit $?
    fi
done


