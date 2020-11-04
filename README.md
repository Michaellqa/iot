# iot

The service aggregates data from multiple data sources and writes the results to a storage of your 
choice.

### Quick start

Execute the commands:
```
docker build -t iot . 
docker run -p 8081:8081 -v $(pwd):/app/data iot
```

The second command mounts the current directory to the docker container so that files created by 
container can be easily viewed.

The app would accepts POST requests on port 8081. First request starts all the components from the 
config. Until the job is done the app will respond with status code 429 (to many requests) to other
requests.
   
```
curl -d "json_config" -X POST http://localhost:8081/
```

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

A generator periodically asks its data sources for the current values and pushes this data into a 
message broker.

#### Message broker

The broker delivers messages from publishers to subscribers. For each subscriber it creates a queue 
that keeps track of subscriber's read offset. Once a message is published it is copied to all the 
queues. A queue has limited capacity. If new messages come to a queue faster then a subscriber 
reads them, and the queue reaches max capacity, the oldest message would be dropped in order to add 
a new one.

#### Aggregators

An aggregator subscribes to the stream of generated data and filters it by id. Selected data is 
saved into a buffer. Once in a period of time aggregator calculates averages in the buffer, saves 
them into a storage and clears the buffer to start a new iteration. Stops aggregating when the data 
channel is closed.

#### Storage

An implementation of aggregation.Storage interface to store aggregated data.
Supported options:

- print to the console
- write to the file
- slow printing, useful for tests and to see how cool async store works. 


