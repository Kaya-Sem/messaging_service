



```mermaid

flowchart LR
  A[Clients]
  B[Gateways<br/>(HTTP API / WebSockets)]
  C[Message Broker<br/>(Kafka / NATS / RabbitMQ)]
  D[Processing Layer<br/>(Consumers / Workers)]
  E[Databases<br/>(Cassandra / Redis / Elasticsearch)]
  F[Real-time Delivery<br/>(WebSocket push / Mobile push / Redis pub/sub)]

  A <-->|connects| B
  B -->|publish message| C
  C -->|consume| D
  D -->|write / archive| E
  D -->|push message| F
  F -->|deliver to gateway| B
  B -->|deliver to client| A

  style A fill:#f3f4f6,stroke:#333,stroke-width:1px
  style B fill:#e0f2fe,stroke:#333,stroke-width:1px
  style C fill:#fef3c7,stroke:#333,stroke-width:1px
  style D fill:#dcfce7,stroke:#333,stroke-width:1px
  style E fill:#ede9fe,stroke:#333,stroke-width:1px
  style F fill:#fbcfe8,stroke:#333,stroke-width:1px


```
