#!/bin/bash

# sysinfo_page - A script to run all comparators

##### Constants

AUTH_TOKEN="1f0d50324de9b11e49e700b31eff7338b7ccf3e196e4f4fd525aa389290237f3"
SCOPE_1="https://read-batch_payment-methods.furyapps.io"
SCOPE_2="https://production-reader_payment-methods-read-v2.furyapps.io"
ARRAY_PATHS=(
  "/Users/mpons/Documents/comparator/payment-methods/v2/LOTE-2/MCO/NONE/MCO-NONE.error"
  "/Users/mpons/Documents/comparator/payment-methods/v2/LOTE-2/MCO/MELI/MCO-MELI.error"
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