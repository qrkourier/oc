FROM centos
RUN cat /etc/resolv.conf
RUN dnf -y install golang make git krb5-devel

# docker build --tag=ocbuilder .
# docker network create ocbuilder
# docker run --rm  --tty --interactive --name=ocbuilder --network=ocbuilder \
# --volume=$PWD:/mnt ocbuilder bash -c "cd /mnt && make oc"

