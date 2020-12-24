# Define some common paths
DATA_OUTPUT_DIR ?= _data
DATA_OUTPUT_FILE ?= ${DATA_OUTPUT_DIR}/data.json
CONFIGS_OUTPUT_DIR ?= _configs


clean:
	rm -r -v -- ${CONFIGS_OUTPUT_DIR}

fetch-data ${DATA_OUTPUT_FILE}:
	go run ./cmd/psxemudatafetch > "${DATA_OUTPUT_FILE}"

update-data: fetch-data
	@(git diff --quiet -- "${DATA_OUTPUT_FILE}" && echo "Data hasn't changed") \
		|| git commit -m "Updating data via fetch" -- "${DATA_OUTPUT_FILE}"

generate-configs ${CONFIGS_OUTPUT_DIR}:
	go run ./cmd/psxemuconf



.PHONY: clean fetch-data update-data generate-configs
