 
# 🚀 Distributed NATS Project (Go + Docker Compose)

This project demonstrates core **NATS** features using Go and Docker Compose, including:

- ✅ High performance and low latency
- ✅ Publish/Subscribe messaging
- ✅ Request/Reply communication
- ✅ Lightweight single-binary client
- ✅ **Distributed NATS cluster** with 3 nodes for horizontal scalability and resilience

---

## 🚦 Getting Started

### 1. Start the full NATS system

```bash
docker compose up -d
````

This launches:

* A **3-node NATS cluster**: `nats1`, `nats2`, `nats3`
* A custom Docker network `natsnet` for name-based service discovery

Each node connects to the others using the `-routes` flag — forming a real **distributed cluster** that can share load and subscriptions.

---

## 🌐 Distributed NATS Architecture

This setup is not just multiple standalone servers — it is a true **distributed system**:

* All nodes share subscription state
* Clients can connect to any node and still receive or send messages across the cluster
* If one node fails, the others continue operating
* The system routes messages across nodes transparently

This distributed architecture provides **horizontal scalability**, **high availability**, and **fast internal routing**.

---

## 💥 Simulate Failover

You can simulate a node failure to test fault tolerance:

```bash
docker compose stop nats1
```

✅ The cluster will remain available using the remaining nodes (`nats2` and `nats3`).
Clients can reconnect or connect to another node without interruption.

If your Go clients are started with:

```go
nats.Connect("nats://nats1:4222,nats://nats2:4222,nats://nats3:4222")
```

...then the client will automatically **failover** to another node when one becomes unreachable.

You can restart the node with:

```bash
docker compose start nats1
```

---

## 📬 Pub/Sub vs 🔁 Req/Rep (Explained)

| Pattern   | Delivery          | Target         | Response? | Use Case           |
| --------- | ----------------- | -------------- | --------- | ------------------ |
| `pub/sub` | Fan-out           | Many listeners | ❌ No      | Broadcasts, events |
| `req/rep` | Load-balanced RPC | One listener   | ✅ Yes     | Queries, RPC calls |

* **Pub/Sub**: Every subscriber receives a copy of the message. No reply is expected.
* **Req/Rep**: Only one replier receives the request and must send back a response.

---

## ⚖️ Kafkaesque vs NATS (Comparison)

| Feature      | Kafka (Kafkaesque)            | NATS                                    |
| ------------ | ----------------------------- | --------------------------------------- |
| Architecture | Log-based                     | Message broker                          |
| Latency      | Higher                        | Very low (sub-ms)                       |
| Storage      | Durable (disk)                | Ephemeral (by default)                  |
| Scaling      | Partition-based               | Clustered, subject-based                |
| Pub/Sub      | Asynchronous fan-out          | Native pub/sub                          |
| Req/Rep      | Not native (extra logic)      | Built-in with `Request()` and `Reply()` |
| Setup        | Heavy (Zookeeper, brokers)    | Lightweight (single binary)             |
| Dev Speed    | Slower startup & ops overhead | Very fast startup, simple setup         |
| Use Case Fit | Event sourcing, analytics     | Realtime messaging, microservices       |

✅ NATS is a better fit for lightweight, high-speed messaging between services.
☁️ Kafka shines in event persistence, replayability, and data streaming.

---

## 🧪 Built-in Client Roles

This Go app supports the following roles:

| Service      | Role  | Description                                 |
| ------------ | ----- | ------------------------------------------- |
| `publisher`  | `pub` | Publishes 100,000 messages to `updates`     |
| `subscriber` | `sub` | Subscribes to `updates` and prints messages |
| `requester`  | `req` | Sends a request to `ping`, waits for reply  |
| `replier`    | `rep` | Replies to `ping` with `pong`               |

Each one runs the same binary with a different `-role` flag.

---

## 🧑‍🔬 Usage Examples

### 🔊 Publish/Subscribe

In two terminals:

```bash
# Terminal A — subscriber listens on subject `updates`
docker compose up subscriber

# Terminal B — publisher sends 100,000 messages
docker compose up publisher
```

---

### 🔁 Request/Reply

In two terminals:

```bash
# Terminal A — replier listens on subject `ping`
docker compose up replier

# Terminal B — requester sends a single message to `ping`
docker compose up requester
```

---

## 🔍 NATS Cluster Details

Your distributed cluster consists of three nodes:

| Service | Port | Cluster Peers |
| ------- | ---- | ------------- |
| nats1   | 4222 | nats2, nats3  |
| nats2   | 4222 | nats1, nats3  |
| nats3   | 4222 | nats1, nats2  |

Each node:

* Listens on port `4222` for client connections
* Uses port `6222` internally to communicate with its peers
* Participates in a **single, unified cluster**

---

## 🧠 Features Demonstrated

| Feature                | Shown In                           |
| ---------------------- | ---------------------------------- |
| High Performance       | `publisher` sends 100k messages    |
| Publish/Subscribe      | `publisher` + `subscriber`         |
| Request/Reply          | `requester` + `replier`            |
| Lightweight Go Binary  | Single compiled app, \~5MB in size |
| Distributed Clustering | 3-node fault-tolerant NATS system  |
| Client Failover        | Tested with `docker compose stop`  |

---

Made with ❤️ to explore distributed messaging with NATS in Go.
