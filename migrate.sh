#!/bin/bash

migrate -path assets/migrations -database sqlite3://data/flimsy.db "$*"
