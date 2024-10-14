
# CartLoom

**CartLoom** is a robust, scalable e-commerce backend system designed to handle high-volume order processing and real-time data synchronization with Shopify. Its goal is to seamlessly integrate with Shopify, process orders efficiently, manage inventory, and provide a reliable platform for handling large-scale transactions in a microservices architecture. 

The system is built to ensure high availability, fault tolerance, and performance, leveraging distributed systems concepts, messaging queues, and global database replication.

## 📋 Objective

The main objective of **CartLoom** is to provide a scalable infrastructure that supports e-commerce platforms like Shopify by integrating them with real-time order processing systems. This includes the following core functionalities:

- **Shopify Integration**: Sync product updates and inventory changes in real-time using Shopify's webhooks.
- **Order Processing with Kafka**: Handle high-volume order processing asynchronously to ensure that large traffic loads do not overwhelm the system.
- **Global Data Availability**: Use Amazon DynamoDB's global tables to ensure data consistency and low-latency access across different geographic regions.
- **Real-Time Monitoring and Metrics**: Monitor the health and performance of the system with Prometheus, and ensure that any issues are detected early with clear and actionable metrics.

This system is designed to be **highly scalable**, **resilient**, and **easy to extend**, making it ideal for large e-commerce stores needing to process thousands of orders per second.

## 🛠️ Technologies Used and Why

### 1. **Go (Golang)**
   - **Why Go?**: Go was chosen for its simplicity, strong performance, and excellent support for concurrency. It is well-suited for building scalable systems that require low-latency and high-throughput operations, making it a perfect match for CartLoom's high-demand requirements.
   - **Benefits**: Fast execution, low resource usage, and a rich standard library for HTTP servers and concurrency management.

### 2. **Kafka**
   - **Why Kafka?**: Kafka is a distributed event streaming platform used for building real-time data pipelines. It was chosen for CartLoom to handle asynchronous, high-throughput order processing. Kafka allows the system to decouple message producers and consumers, ensuring smooth data flow even during peak traffic times.
   - **Benefits**: High fault tolerance, scalability, and the ability to process millions of messages per second. It’s ideal for applications where message ordering and durability are critical.

### 3. **Redis**
   - **Why Redis?**: Redis is used as a distributed in-memory cache to speed up data access and reduce the load on the primary database. For CartLoom, Redis is crucial for reducing latencies in reading and writing frequently accessed data like product inventory and session data.
   - **Benefits**: Extremely fast read/write operations, simple to use, and supports various data structures like strings, hashes, and lists. It also supports distributed locks, which is key for preventing race conditions in inventory management.

### 4. **Amazon DynamoDB**
   - **Why DynamoDB?**: DynamoDB is a fully managed NoSQL database with built-in support for global tables, making it the perfect solution for CartLoom’s need to replicate data across multiple regions for high availability and low-latency reads. The choice of DynamoDB ensures that the system can handle millions of requests per second without downtime or data loss.
   - **Benefits**: Global replication, serverless scalability, and automatic backup. DynamoDB's pay-per-request model makes it cost-effective for both small and large-scale operations.

### 5. **Shopify API**
   - **Why Shopify API?**: CartLoom integrates with Shopify to handle product updates, orders, and inventory in real-time. The Shopify API provides a powerful interface for e-commerce stores to interact with their data programmatically.
   - **Benefits**: Seamless integration with Shopify stores, access to product and order data, and real-time synchronization through webhooks.

### 6. **Prometheus & Grafana**
   - **Why Prometheus?**: Prometheus is used to monitor the health of CartLoom’s services. Metrics collected from each service provide visibility into system performance and help identify bottlenecks or failures. Prometheus also integrates with Grafana for detailed, real-time data visualization.
   - **Benefits**: Real-time metrics, efficient time-series storage, and robust alerting features for monitoring the system.
   - **Why Grafana?**: Provides beautiful, easy-to-understand dashboards for visualizing the data collected by Prometheus.
   
### 7. **Docker & Kubernetes**
   - **Why Docker?**: Docker is used to package and run CartLoom’s services in isolated containers, ensuring consistency across different environments (development, testing, production).
   - **Benefits**: Simplifies dependency management, ensures consistency, and isolates each service in its own container.
   - **Why Kubernetes?**: Kubernetes is used for orchestrating and scaling the microservices that make up CartLoom. It automates deployment, scaling, and operation of application containers, allowing CartLoom to scale horizontally based on demand.
   - **Benefits**: Autoscaling, self-healing, and easy management of containerized applications at scale.

### 8. **Istio & Consul**
   - **Why Istio?**: Istio is used as a service mesh to manage communication between microservices, providing security (mTLS), traffic management, and observability.
   - **Why Consul?**: Consul provides service discovery and distributed configuration, allowing CartLoom to manage configurations and ensure services are dynamically detected and updated as needed.
   - **Benefits**: Better control over microservice communications, load balancing, secure service-to-service communication, and service discovery.

## 💻 How to Run the Project

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/cartloom.git
cd cartloom
```

### 2. Set up Environment Variables

Create a `.env` file based on the provided **`env.example`** file:

```bash
cp .env.example .env
```

Edit the `.env` file and add your own configurations for Redis, DynamoDB, and Shopify API credentials.

### 3. Build and Run the Application

Install Go dependencies:

```bash
go mod tidy
```

Start the services (Kafka, Redis, Prometheus, etc.) using Docker Compose:

```bash
docker-compose up --build
```

Run the Go application:

```bash
go run main.go
```

## System Architecture

The system is designed using a **microservices architecture** that leverages distributed systems principles. Here's a high-level overview of the components:

- **Shopify**: Sends product updates via webhooks.
- **Kafka**: Handles real-time order processing and communication between services.
- **Redis**: Used as a caching layer for frequently accessed data.
- **DynamoDB**: Stores orders and product data, with global replication for high availability.
- **Prometheus & Grafana**: Used to monitor system performance and health.
- **Docker & Kubernetes**: Container orchestration and management of microservices.
