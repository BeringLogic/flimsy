#!/bin/bash

/home/phil/.local/bin/migrate -path src/assets/migrations -database sqlite3://data/flimsy.db $*
