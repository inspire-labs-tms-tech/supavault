#!/bin/bash

# globals
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# path to Version.java relative to the script's directory
VERSION_FILE="$SCRIPT_DIR/src/main/java/com/inspiretmstech/supavault/constants/Version.java"

# test app version
if [ -z "${APP_VERSION}" ]; then
    echo "'APP_VERSION' is not defined in the runtime environment"
    exit 1
fi

# update the app version
# sed is different on mac than linux
# => need to handle appropriately for local versus cloud build
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/\"0.0.0\"/\"${APP_VERSION}\"/" "${VERSION_FILE}"
else
    # Linux and other Unix-like systems
    sed -i "s/\"0.0.0\"/\"${APP_VERSION}\"/" "${VERSION_FILE}"
fi

# build the fat jar
gradle shadowJar

# cleanup any old pkg files
sudo rm -rf ./supavault-*.pkg

# build new pkg
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    jpackage \
      --verbose \
      --name supavault \
      --main-jar supavault.jar \
      -i ./build/libs \
      --install-dir /usr/local \
      -t pkg \
      --resource-dir ./scripts \
      --app-version "${APP_VERSION}"
else
    # Linux and other Unix-like systems
    jpackage \
      --verbose \
      --name inspire-tms \
      --main-jar main.jar \
      -i ./main/build/libs \
      --linux-package-name inspire-tms
fi
