FROM scratch
EXPOSE 8080
ENTRYPOINT ["/book-microservice"]
COPY ./bin/ /