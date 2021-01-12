#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="db9bb6a58125e2171ef556fe808a0ffec6838dc725aac79244514c16cc924399"
SCOPE_1="http://api.mp.internal.ml.com/test"
#SCOPE_2="https://production-reader_payment-methods-read-v2.furyapps.io"
SCOPE_2="http://api.mp.internal.ml.com"
ARRAY_PATHS=(
  "/Users/mpons/Downloads/bin-api/BINS_ALL_SITES_100000-aa.csv"

  #"/Users/mpons/Documents/comparator/payment-methods/v2/LOTE-3/MCO-LOTE-3.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/MLM.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/NONE/MLM/MLM.error"

	)

#ARRAY_CHANNELS=("" "point" "splitter" "instore")
ARRAY_CHANNELS=("")

for i in "${ARRAY_PATHS[@]}"
do
	for j in "${ARRAY_CHANNELS[@]}"
	do
		jsonparator -path "$i" -host "${SCOPE_1}" -host "${SCOPE_2}" -V 1 -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -E "settings.#.version"
	done
done