FROM centurylink/ca-certs

ADD drone-discord /

ENTRYPOINT ["/drone-discord"]
