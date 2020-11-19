#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="764b1baf6ddbff24d29ac2b7d70483f7b4768701c667549a87ca98aeb0492fd6"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
#SCOPE_2="https://production-reader_payment-methods-read-v2.furyapps.io"
SCOPE_2="https://testing-reader-data_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
  "/Users/mpons/Downloads/test-dieguin.csv"

  #"/Users/mpons/Documents/comparator/payment-methods/v2/LOTE-3/MCO-LOTE-3.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/MELI/MLM/MLM.error"
	#"/Users/mpons/Documents/comparator/payment-methods/v2/5_17-08-2020_21-08-2020/202008-10-15/NONE/MLM/MLM.error"

	)

#ARRAY_CHANNELS=("" "point" "splitter" "instore")
ARRAY_CHANNELS=("point")

for i in "${ARRAY_PATHS[@]}"
do
	for j in "${ARRAY_CHANNELS[@]}"
	do
		jsonparator -path "$i" -host "${SCOPE_1}" -host "${SCOPE_2}" -V 5 -header "X-Auth-Token:${AUTH_TOKEN}" -header "X-Caller-Scopes:$j" -E "paging" -E "results.#.payer_costs.#.payment_method_option_id"
	done
done