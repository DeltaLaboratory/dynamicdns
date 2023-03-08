# Self-Hosted Dynamic DNS (DDNS) Service

This is a cool and easy-to-use self-hosted DDNS service that allows you to assign a domain to your dynamic IP address. It's perfect for home networks, small businesses, or anyone who needs to access their devices remotely without having to remember their IP address.

## Configuration
We're using [HCL](https://github.com/hashicorp/hcl) for the configuration file format.
```hcl
interval = 1

ddns {
  service "cloudflare" {
    api_key = "API_KEY"
    zone_id = "ZONE_ID"
    record_name = "example.example.com"
    ttl = -1
  }
}
```
* "interval" refers to the frequency, measured in minutes, at which the record should be updated.
* The "cloudflare" service refers to the DNS provider for the zone, and currently, only Cloudflare is supported.
* "api_key" refers to the API key for the DNS provider
* "zone_id" refers to the zone ID of the DNS record.
* "ttl" stands for Time-to-Live and specifies how long the record should be cached, with a value of -1 indicating that it should be set automatically.

## Installation
### With Docker Compose
```yaml
services:
  ddns:
    image: ghcr.io/DeltaLaboratory/dynamicdns:latest
    container_name: ddns
    restart: unless-stopped
    networks:
      - proxy
    volumes:
      - ./config/:/ko-app/
```
also do not forget to create config file in ./config/config.hcl
### Build yourself
* require go 1.19 or later (do not guarantee to work older version)
```shell
git clone github.com/DeltaLaboratory/dynamicdns
cd dynamicdns
go mod download
go build -o ddns
./ddns
```
also do not forget to create config file in ./config.hcl

## Usage

Once the service is running, it will automatically update your DNS record every time your IP address changes. You can access your device or server using the hostname you specified in the `config.hcl` file.


If you're still having issues, feel free to open an issue on this repository.

## Contributing

Contributions to this project are welcome! If you find a bug or have a feature request, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
