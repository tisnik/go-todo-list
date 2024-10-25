#!/bin/sh
 
DATABASE=todo.db
 
cat "db_init.sql" | sqlite3 "${DATABASE}"
