# Tiny Docker image for Go.
FROM iron/go

# Copy solar binary
COPY solar-grpc-controller /

# Set entrypoint to the solar binary
ENTRYPOINT ["/solar-grpc-controller"]
