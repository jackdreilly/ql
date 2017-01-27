FROM golang:onbuild
ADD static/components static/components
ADD static/css static/css
ADD static/quik_lyrics.html static/quik_lyrics.html
ADD static/not_found.html static/not_found.html
ADD static/bower.json static/bower.json
RUN     apt-get update
RUN     apt-get install -y node npm
WORKDIR /app/static/
RUN     npm install -g bower
RUN     bower install --allow-root
EXPOSE 8080