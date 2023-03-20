FROM redhat/ubi8

# Install podman
RUN yum install -y \
      podman \
      git \
      golang \
      make

# Install docker
# Based on guide here:
# https://faun.pub/how-to-install-simultaneously-docker-and-podman-on-rhel-8-centos-8-cb67412f321e
RUN yum install -y \
      yum-utils && \
    yum-config-manager \
      --add-repo \
      https://download.docker.com/linux/centos/docker-ce.repo && \
    # We only need to install the CLI to use the Docker socket shared from the
    # host, which is enough for pulling information from Docker.
    yum install -y docker-ce-cli

# Install FIO from source, since it's not available from the
# UBI package repositories unless also running on a registered
# RHEL host.
RUN yum install -y \
    wget && \
    #
    #  Build libaio for iops tests
    #
    wget \
      -O libaio.tar.gz \
      https://pagure.io/libaio/archive/libaio-0.3.111/libaio-libaio-0.3.111.tar.gz && \
    mkdir libaio && \
    tar \
      --strip-components=1 \
      --directory=libaio \
      -xvf libaio.tar.gz && \
    cd libaio && \
    make install && \
    ldconfig && \
    cd .. && \
    rm libaio.tar.gz && \
    #
    # Build fio
    #
    wget \
      -O fio.tar.gz \
      https://github.com/axboe/fio/archive/refs/tags/fio-3.33.tar.gz && \
    mkdir fio && \
    tar \
      --strip-components=1 \
      --directory=fio \
      -xvf fio.tar.gz && \
    cd fio && \
    ./configure && \
    make && \
    make install && \
    cd .. && \
    rm -rf fio.tar.gz && \
    yum remove -y wget
