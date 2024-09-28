gradle build
jpackage \
  --verbose \
  --name supavault \
  --main-jar cli.jar \
  --main-class com.inspiretmstech.supavault.Main \
  -i ./cli/build/libs \
  --install-dir /usr/local \
  -t pkg \
  --resource-dir ./scripts
