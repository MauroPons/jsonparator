#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="fff0f1565a16180d42002def5e1c91fd733c9d81170addee08a4ec1154e55bb8"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
SCOPE_2="https://production-reader-testscope_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
  "/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/NONE/MCO/MCO.error"
  #"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/MELI/MCO/MCO.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/1_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/MLM.csv"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/1_17-08-2020_21-08-2020/202008-10-15/NONE/MLM/MLM.csv"

	)

#ARRAY_CHANNELS=("" "point" "splitter" "instore")
ARRAY_CHANNELS=("")

for i in "${ARRAY_PATHS[@]}"
do
	for j in "${ARRAY_CHANNELS[@]}"
	do
		jsonparator -path "$i" -host "${SCOPE_1}" -host "${SCOPE_2}" -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -exclude "paging" -exclude "results.#.payer_costs.#.payment_method_option_id" -M "marketplace"
	done
done