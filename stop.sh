#!/bin/bash

# Set script to exit on error
set -e

# Define the database file
DB_FILE="task.db"


# Delete the database if it exists
if [ -f "$DB_FILE" ]; then
  echo "Deleting database..."
  rm -f "$DB_FILE"
fi

echo "database deleted."
