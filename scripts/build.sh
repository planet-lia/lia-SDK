#!/usr/bin/env bash

runTests=$1

platforms=("linux/amd64" "linux/386" "windows/amd64" "windows/386" "darwin/amd64" "linux/amd64/static")

# Cd to root of the project
pathToScript="`dirname \"$0\"`"
cd ${pathToScript}/..

# Run tests
if [[ $runTests != "false" ]]; then
    go test ./...

    exit_status=$?
    if [[ ${exit_status} != 0 ]]; then
        echo ${exit_status}
        (>&2 echo "Running tests failed.")
        exit ${exit_status}
    fi
fi

for platform in "${platforms[@]}"
do
    echo "Building for ${platform}..."

    platformSplit=(${platform//\// })
    GOOS=${platformSplit[0]}
    GOARCH=${platformSplit[1]}


    if [[ ${GOARCH} == "amd64" ]]; then
        archName="x64"
    else
        archName="x32"
    fi

    buildDir="build/lia-sdk-"${GOOS}'-'${archName}

    if [[ ${#platformSplit[@]} == 3 ]]; then
       CGO_ENABLED=0
       buildDir=${buildDir}-static
    fi

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

    if [[ ${CGO_ENABLED} == 0 ]]; then
        env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${buildDir}/${execName} cmd/lia/main.go
    else
        env GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${buildDir}/${execName} cmd/lia/main.go
    fi
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit $?
    fi
done


