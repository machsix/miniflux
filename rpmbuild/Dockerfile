FROM centos:latest
RUN yum install -y rpm-build
RUN yum install -y redhat-rpm-config
RUN mkdir -p /root/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
RUN echo "%_topdir /root/rpmbuild" >> .rpmmacros
