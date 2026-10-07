### Multi-platform

To be run from the root folder.

`docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile -t filestash:buildx --load .`
