# Define some common paths
DATA_OUTPUT_DIR ?= _data
DATA_OUTPUT_FILE ?= ${DATA_OUTPUT_DIR}/data.json


fetch-data ${DATA_OUTPUT_FILE}:
	go run ./cmd/psxemudatafetch > "${DATA_OUTPUT_FILE}"

update-data: fetch-data
	git commit -m "Updating data via fetch" -- "${DATA_OUTPUT_FILE}"



.PHONY: fetch-data update-data
