#!/bin/bash
goose -dir migrations create $1 sql
