#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="3c67933db5eb668049b1926f71656ead95e36b3faa0efc2ec16186a262b0d11f"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
SCOPE_2="https://production-reader-testscope_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
	"/Users/mpons/Documents/comparator/payment-methods/v2/3_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/total.error"
	"/Users/mpons/Documents/comparator/payment-methods/v2/3_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/total.error"
	"/Users/mpons/Documents/comparator/payment-methods/v2/3_17-08-2020_21-08-2020/202008-10-15/NONE/MCO/total.error"
	"/Users/mpons/Documents/comparator/payment-methods/v2/3_17-08-2020_21-08-2020/202008-10-15/NONE/MCO/total.error"
	)

#ARRAY_CHANNELS=("" "point" "splitter" "instore")
ARRAY_CHANNELS=("")

for i in "${ARRAY_PATHS[@]}"
do
	for j in "${ARRAY_CHANNELS[@]}"
	do
		jsonparator -path "$i" -host "${SCOPE_1}" -host "${SCOPE_2}" -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -exclude "paging" -exclude "results.#.payer_costs.#.payment_method_option_id" -M "marketplace" -S true
	done
done