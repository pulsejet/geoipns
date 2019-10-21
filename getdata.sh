#!/bin/bash

# Move to data directory
mkdir -p ./data
cd data

# Download
wget "https://geolite.maxmind.com/download/geoip/database/GeoLite2-City-CSV.zip"
wget "https://geolite.maxmind.com/download/geoip/database/GeoLite2-ASN-CSV.zip"

# Cleanup old
rm -rf GeoLite2-City-CSV_*
rm -rf GeoLite2-ASN-CSV_*

# Unzip
unzip GeoLite2-City-CSV.zip
unzip GeoLite2-ASN-CSV.zip

# Get files
rm -rf GeoLite*.csv
mv ./GeoLite2-City-CSV_*/*.csv .
mv ./GeoLite2-ASN-CSV_*/*.csv .

# Cleanup
rm -rf GeoLite2-City-CSV_*
rm -rf GeoLite2-ASN-CSV_*
rm -rf GeoLite2*.zip

cd ..
