# toll-calculator

## Installing protobuf compilers

For linux users or WSL2
```
sudo apt install -y protobuf-compiler
```

For mac users use brew for this

```
brew install protobuff
```

## Installing GRPC and Protobuff plugins for golang
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```
2. GRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

3. Note: you need to set the go/bin directory into your PATH
export PATH="${PATH}:${HOME}/go/bin"

4.  Install package dependencies

4.1 Protobuffer package
```
go get google.golang.org/protobuf/
```

4.2 grpc package
```
go get google.golang/grpc/
```

```
docker run -d --name kafka --hostname kafka -p 9092:9092 \
-e KAFKA_ENABLE_KRAFT=yes \
-e KAFKA_CFG_NODE_ID=1 \
-e KAFKA_CFG_BROKER_ID=1 \
-e KAFKA_CFG_PROCESS_ROLES=controller,broker \
-e KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
-e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
-e KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
-e KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@:9093 \
-e KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER \
bitnami/kafka:latest

```

```
docker run -d \
--rm \
--name zookeeper-1 \
--net kafka \
-v ${PWD}/config/zookeeper-1/zookeeper.properties:/kafka/config/zookeeper.properties \
okusa/zookeeper:2.7.0
```
```
docker run -d \
--rm \
--name kafka-1 \
--net kafka s\
-v ${PWD}/config/kafka-1/server.properties:/kafka/config/server.properties \
okusa/kafka:2.7.0


docker logs kafka-1
```

```
docker run -d \
--rm \
--name kafka-2 \
--net kafka \
-v ${PWD}/config/kafka-2/server.properties:/kafka/config/server.properties \
okusa/kafka:2.7.0

```

```
docker run -d \
--rm \
--name kafka-3 \
--net kafka \
-v ${PWD}/config/kafka-3/server.properties:/kafka/config/server.properties \
okusa/kafka:2.7.0

```

```
/kafka/bin/kafka-topics.sh \
--create \
--zookeeper zookeeper-1:2181 \
--replication-factor 1 \
--partitions 3 \
--topic Orders

```

```

/kafka/bin/kafka-topics.sh \
--describe \
--topic Orders \
--zookeeper zookeeper-1:2181

```

```
Topic: Orders   PartitionCount: 3       ReplicationFactor: 1    Configs: 
Topic: Orders   Partition: 0    Leader: 3       Replicas: 3 Isr: 3
Topic: Orders   Partition: 1    Leader: 1       Replicas: 1 Isr: 1
Topic: Orders   Partition: 2    Leader: 2       Replicas: 2 Isr: 2
```