# Fios Exporter
This is a Prometheus Exporter for the Verizon Fios Quantum Gateway Router that I made for my Grafana Dashboard. Currently it only does bandwidth stats, specifically the average kilobits/second over the last minute, which is the closet you'll get to what's currently passing through your router. I'm happy to add more stats if people are interested, feel free to open an issue or even a pull request.

## Metrics

`fios_tx_kbs_minute_1`: The average kilobits/second of outgoing bandwidth over the last minute

`fios_rx_kbs_minute_1`: The average kilobits/second of incoming bandwidth over the last minute

## Setup

### Flags:

`--password`: The password to the admin dashboard of your router (required).

`--host`: The address to your router (optional). Default: https://myfiosgateway.com

`--port`: The port where the metrics are hosted (optional). Default: 2100

### Docker (recomended):

CLI:

`docker run -d -p 2100:2100 ghcr.io/leviharrison/fios-exporter --password=<you_password>`

Docker compose:

```
version: '3'
  services:
    fios-exporter:
      image: ghcr.io/leviharrison/fios-exporter
      command:
        - "--password=<your_password>"
      ports:
        - "2100:2100"
```
