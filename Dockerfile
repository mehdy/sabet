FROM ubuntu:latest

ADD ./sabet /usr/local/bin/sabet

CMD [ "./sabet" ]