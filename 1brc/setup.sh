#!/usr/bin/env sh

wget \
	-O - https://raw.githubusercontent.com/gunnarmorling/1brc/main/data/weather_stations.csv \
	> ./data/weather_stations.csv \
&&
	python3 ./data/create_measurements.py $1
