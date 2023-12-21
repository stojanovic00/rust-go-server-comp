# Metrics

- response time
- throughput
- CPU utilization
- Mem utilization
- Disk utilization
- network utilization


# Tools:

- jmeter

# Requests

- read only (GET)
    - tests read lock on map and entry level

- write only (PUT with new entities)
    - tests write lock on map level (also read lock)

- read and write (PUT on existing entities, each entity 2 times)
    - 1st tests read lock on map
    - 2nd tests write lock on an entry
