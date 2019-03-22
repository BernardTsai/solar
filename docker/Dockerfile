# Tiny Docker image for Go.
FROM iron/go

# Copy solar binary
COPY solar /

# Copy solar configuration file
COPY solar-conf.yaml /

# Copy static web content
COPY static/ /static

# Set entrypoint to the solar binary
ENTRYPOINT ["/solar"]
