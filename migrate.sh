#!/bin/bash

/home/phil/.local/bin/migrate -path assets/migrations -database sqlite3://data/flimsy.db $*
