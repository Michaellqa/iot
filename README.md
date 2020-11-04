# iot

Service aggregates data from multiple data sources and writes it into a single storage.

Has a web API, listens for requests on port 8081.

### Quick start

```
docker build -t iot . 
docker run -v $(pwd):/app/data iot

curl -d "json_config" -X POST http://localhost:8081/
```

The second command mounts your current directory to the docker container so that files created by container with can be viewed.

Config example:

```json
{
  "generators": [
    {
      "timeout_s": 20,
      "send_period_s": 1,
      "data_sources": [
        {
          "id": "data_1",
          "init_value": 50,
          "max_change_step": 5
        }
      ]
    }
  ],
  "aggregators": [
    {
      "aggregation_period_s": 5,
      "sub_ids": [
        "data_1"
      ]
    }
  ],
  "queue": {
    "size": 50
  },
  "storage_type": 0
}
``` 

### Components

#### Generators

A generator periodically asks its data sources for the current values and push this data into a message broker.

#### Message broker

Broker delivers messages from publishers to subscribers. 
For each subscriber it creates a queue that keeps track of subscriber's read offset. 
Once a message is published it is copied to all the queues. A queue has limited capacity.
If new messages come to the queue faster then the subscriber reads them, and the queue reaches 
max capacity, the oldest message would be dropped in order to add a new one.

#### Aggregators

An aggregator subscribe to the stream of data and filters them by id. Targeted data is saved into a 
buffer. Once in a period of time aggregator calculates averages in the buffer, saves them into a 
storage and clears the buffer to start a new iteration.


