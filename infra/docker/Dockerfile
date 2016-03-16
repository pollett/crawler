FROM golang

RUN apt-get update && apt-get -y install inotify-tools psmisc
RUN git clone https://github.com/alexedwards/go-reload.git

RUN mv go-reload/go-reload /usr/local/bin/
RUN chmod +x /usr/local/bin/go-reload

ENV GOPATH=/go/src/app/vendor:/go/src/app
