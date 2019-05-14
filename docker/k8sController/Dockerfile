# Tiny Docker image for Go.
FROM iron/go

# Copy k8s controller binary
COPY solar-k8s-controller /

# Set entrypoint to the k8s controller binary
ENTRYPOINT ["/solar-k8s-controller"]
