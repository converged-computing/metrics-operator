FROM ubuntu

# docker build -t ghcr.io/converged-computing/benchmark-sysstat .
# The entrypoint to download the pid wait command will be generated on the fly
# via a config map

RUN apt-get update && apt-get install -y git g++ build-essential
RUN git clone https://github.com/sysstat/sysstat /opt/sysstat && \
    cd /opt/sysstat && \
    ./configure --prefix=/usr && \
    make && \
    make install