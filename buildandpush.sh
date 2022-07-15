if [ -z ${platform+x} ]
then
  echo docker build --tag=registry.pi.local:21212/koe:latest -f=Dockerfile .
else
  docker build --tag=registry.pi.local:21212/koe:latest --platform ${platform} -f=Dockerfile .
fi

docker push registry.pi.local:21212/koe:latest
