#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="edc3cad5bd9797883138b2e0ea00f676d2f066a4588ebc5d02ba72333afa8b40"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
#SCOPE_2="https://production-reader_payment-methods-read-v2.furyapps.io"
SCOPE_2="https://testing-reader-comparator_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
  "/Users/mpons/Downloads/MCSPM-2916-400-400.txt"

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
		jsonparator -path "$i" -host "${SCOPE_1}" -host "${SCOPE_2}" -V 5 -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -E "paging" -E "results.#.payer_costs.#.payment_method_option_id"
	done
done