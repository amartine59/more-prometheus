# Prometheus Playground

Trying out Prometheus instrumentation for a small Go app (Visualizing metrics with Grafana).
This is a very simple example but it should allow you to play a bit with both prometheus and Grafana and see what they have to offer.

## Requirements

To run this one you only need docker installed.
Additionally, if you wish to generate some HTTP calls against the root of the web app, you can use
[Hey](https://github.com/rakyll/hey) and run the following command:

```
hey -z 5m -q 5 -m GET -H "Accept: text/html" http://127.0.0.1:2112
```

## Running the stack
You can run the stack of services simply by running:
```
docker-compose up -d
```
`Prometheus` should be available at `localhost:9090`.
`Grafana` should be reachable at `localhost:3000`.(u:`admin` - pw: `admin`)

### References
This repository was created following the instructions provided by this nice [Blog Post](https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/)

