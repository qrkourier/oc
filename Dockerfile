FROM centos:centos8
ARG go_distribution_file=go1.16.9.linux-amd64.tar.gz
RUN dnf -y install make git krb5-devel which wget
RUN dnf -y group install "Development Tools"
RUN wget -q https://golang.org/dl/${go_distribution_file}
RUN tar -xzf ${go_distribution_file} -C /usr/local/
RUN mkdir /root/go
ENV GOPATH=/root/go
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH 
WORKDIR /mnt
RUN go env
CMD make oc

# docker build --tag=ocbuilder .

# docker run --rm  --tty --interactive --name=ocbuilder --network=ocbuilder --volume=$PWD:/mnt ocbuilder 

# ./oc version
