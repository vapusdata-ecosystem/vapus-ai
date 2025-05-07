package setupconfig

var TrinoSetup = `
fullnameOverride: vapusdata-trino
coordinatorNameOverride: vapusdata-t-coordinator
workerNameOverride: vapusdata-t-worker
service:
    port: 8088
coordinator:
    resources:
        requests:
            cpu: 500m
            memory: 1Gi
    jvm:
        maxHeapSize: 2G
        gc: G1GC
worker:
    replicas: 2
    resources:
        requests:
            cpu: 500m
            memory: 1Gi
config:
    coordinator:
        logLevel: INFO
    properties:
        query.max-memory: 8GB
        query.max-memory-per-node: 2GB
`
