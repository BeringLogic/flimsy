#!/bin/bash

migrate -path src/assets/migrations -database sqlite3://data/flimsy.db $*
