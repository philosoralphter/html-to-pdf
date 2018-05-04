FROM centos

#pull in chrome binary
#ADD ./chrome/chromedriver /usr/local/bin/chrome

#Create Chrome yum rpm then install
RUN /bin/bash -c "echo $'[google-chrome]\n\
name=google-chrome - \$basearch\n\
baseurl=http://dl.google.com/linux/chrome/rpm/stable/\$basearch\n\
enabled=1\n\
gpgcheck=1\n\
gpgkey=https://dl-ssl.google.com/linux/linux_signing_key.pub' > /etc/yum.repos.d/google-chrome.repo"

RUN yum install -y epel-release
RUN yum install -y libappindicator-gtk3
RUN yum install -y google-chrome-stable
RUN yum install -y google-noto-sans-fonts
RUN yum install -y google-noto-serif-fonts

ENV CHROME_LOCATION=/opt/google/chrome/chrome

EXPOSE 9190/tcp


# Now just add the binary
ADD ./output/html-to-pdf /usr/local/bin/html-to-pdf-server

CMD ["/usr/local/bin/html-to-pdf-server"]
