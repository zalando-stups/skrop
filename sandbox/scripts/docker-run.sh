#!/usr/bin/env bash
set -e

images=$(docker images -f "label=artifactId=skrop-sandbox" -q)
if [ ! -z "$images" ]; then
  docker run -it -e NODE_ENV=production --name skrop-sandbox -d -p 3000:3000 ${images}

  echo 'open in browser: http://localhost:3000'
else
    echo "Image to run not found. Run ./docker-build.sh first."
fi
