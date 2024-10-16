gradle shadowJar
sudo rm -rf supavault-*.pkg
jpackage \
  --verbose \
  --name supavault \
  --main-jar supavault.jar \
  -i ./build/libs \
  --install-dir /usr/local \
  -t pkg \
  --resource-dir ./scripts
