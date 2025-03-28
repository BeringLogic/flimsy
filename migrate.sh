#!/bin/bash

migrate -path src/migrations -database sqlite3://data/flimsy.db $*
