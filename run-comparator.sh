#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="c6fa9d383a9998d55f09e67920e293645b03ae9eea68a16d613e1e547fea33a3"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
SCOPE_2="https://production-reader-testscope_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
  "/Users/mpons/Documents/comparator/payment-methods/v2/1_17-08-2020_21-08-2020/202008-10-15/NONE/MCO/MCO.csv"
  "/Users/mpons/Documents/comparator/payment-methods/v2/1_17-08-2020_21-08-2020/202008-10-15/MELI/MCO/MCO.csv"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/MLM.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/NONE/MLM/MLM.error"

	)

#ARRAY_CHANNELS=("" "point" "splitter" "instore")
ARRAY_CHANNELS=("")

for i in "${ARRAY_PATHS[@]}"
do
	for j in "${ARRAY_CHANNELS[@]}"
	do
		jsonparator -path "$i" -host "${SCOPE_1}" -V 5 -host "${SCOPE_2}" -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -E "paging" -E "results.#.payer_costs.#.payment_method_option_id" -M "marketplace"
	done
done