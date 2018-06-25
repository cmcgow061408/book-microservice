FROM scratch
EXPOSE 8080
ENTRYPOINT ["/cmcgow"]
COPY ./bin/ /