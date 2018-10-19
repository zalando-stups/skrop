## Upgrade dependencies

- Update the alpine version
- Update the skipper version
- Update the bimg version
- Update the golang version
- Update the libvips version

```bash
brew upgrade vips
brew upgrade glide
brew upgrade golang 
glide update

# to clear the cache in case glide update failes
glide cc

glide install

make build-docker-vips
```
