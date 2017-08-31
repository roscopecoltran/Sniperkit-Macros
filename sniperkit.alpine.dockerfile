###########################################################################
#		  
#  Build the image:                                               		  
#    $ docker build -t sniperkit -f sniperkit.alpine.dockerfile --no-cache . 		# longer but more accurate
#    $ docker build -t sniperkit -f sniperkit.alpine.dockerfile . 					# faster but increase mistakes
#                                                                 		  
#  Run the container:                                             		  
#    $ docker run -it --rm -v $(pwd)/shared:/shared -p 4242:4242 sniperkit
#    $ docker run -d --name sniperkit -p 4242:4242 -v $(pwd)/shared:/shared sniperkit
#                                                              		  
###########################################################################

## LEVEL1 ###############################################################################################################

FROM alpine:3.6
LABEL maintainer "Luc Michalski <michalski.luc@gmail.com>"

ARG GOSU_VERSION=${GOSU_VERSION-="1.10"}

# Install Gosu to /usr/local/bin/gosu
ADD https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-amd64 /usr/local/sbin/gosu

# Install runtime dependencies & create runtime user
RUN chmod +x /usr/local/sbin/gosu \
 && apk --no-cache --no-progress add ca-certificates git libssh2 openssl \
 && adduser -D app -h /data -s /bin/sh

# Copy source code to the container & build it
COPY . /app
WORKDIR /app
RUN ./templates/docker/sniperkit/scripts/build.sh

# NSSwitch configuration file
COPY ./templates/docker/sniperkit/conf/nsswitch.conf /etc/nsswitch.conf

# App configuration
ENV SKIT_G2E_REPO_PATH "/data/repo"

# Container configuration
VOLUME ["/data"]
EXPOSE 4242
CMD ["/usr/local/sbin/gosu", "app", "/app/sniperkit"]