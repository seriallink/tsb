psql -U postgres -d homework -f data/cpu_usage.sql
psql -U postgres -d homework -c "\COPY cpu_usage FROM data/cpu_usage.csv CSV HEADER"
