#!/bin/bash

migrate -path src/internal/db/migrations -database sqlite3://data/flimsy.db $*
