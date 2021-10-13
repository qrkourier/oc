FROM centos
ARG go_distribution_file=go1.16.9.linux-amd64.tar.gz
RUN cat /etc/resolv.conf
RUN dnf -y install make git krb5-devel which
RUN dnf -y install wget
RUN wget -q https://golang.org/dl/${go_distribution_file}
RUN tar -xzf ${go_distribution_file} -C /usr/local/
ENV GOROOT=/usr/local/go
ENV GOPATH=/mnt
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH 
WORKDIR /mnt
RUN go env
RUN go mod tidy
RUN go mod vendor
RUN make oc

# docker build --tag=ocbuilder .

