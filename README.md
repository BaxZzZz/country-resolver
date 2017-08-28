[![Build Status](https://travis-ci.org/BaxZzZz/country-resolver.svg?branch=master)](https://travis-ci.org/BaxZzZz/country-resolver) 

# country-resolver

Description:

GeoIP information aggregator service with local cache. Demo.

Application launch:

In order to run the application install >=MongoDB 3.0.14. and package manager:

<pre>
go get github.com/constabulary/gb/...
</pre>

To build the application run the following command:

<pre>
gb build
</pre>

To launch unit tests run the command: 

<pre>
gb test
</pre>

Configuration:

In order to get the default configuration file launch the application to generate it.Configuration file:

<pre>
{
    "tcp_server": {
        "address": "0.0.0.0:9999" #
    },
    "geo_ip_provider": {
        "providers": [
            "freegeoip.net",
            "geoip.nekudo.com"
        ],
        "requests_limit": 100,
        "time_interval_min": 1
    },
    "cache": {
        "address": "localhost",
        "db_name": "resolver",
        "username": "",
        "password": "",
        "collection": "cache",
        "items_limit": 100000
    }
}

</pre>

Application check:

To check the work of the application one can use Netcat or Telnet, for example:

<pre>
nc 0.0.0.0 9999
</pre>

The service must return the information about the country based on the user's IP address.
