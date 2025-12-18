version: '3.9'

networks:
  ghost_network:
    driver: bridge

services:
{{- range .Services }}
  # ==========================================
  # App: {{ .Name }}
  # ==========================================
  {{ .Name }}:
    build: {{ .Path }}
    container_name: {{ .Name }}
    networks:
      - ghost_network
    ports:
      - "{{ .Port }}:{{ .Port }}"
    restart: unless-stopped
    environment:
      - SERVER_PORT={{ .Port }}
      - SPRING_PROFILES_ACTIVE=docker
      - MANAGEMENT_ENDPOINTS_WEB_EXPOSURE_INCLUDE=*
      - MANAGEMENT_ENDPOINT_HEALTH_SHOWDETAILS=always

      # --- Databases ---
      {{- if .Features.HasPostgres }}
      - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/{{ $.DbName }}
      - SPRING_DATASOURCE_USERNAME={{ $.DbUser }}
      - SPRING_DATASOURCE_PASSWORD={{ $.DbPass }}
      - SPRING_JPA_HIBERNATE_DDL_AUTO=update
      {{- end }}

      {{- if .Features.HasMySQL }}
      - SPRING_DATASOURCE_URL=jdbc:mysql://mysql:3306/{{ $.DbName }}?useSSL=false&allowPublicKeyRetrieval=true&createDatabaseIfNotExist=true
      - SPRING_DATASOURCE_USERNAME={{ $.DbUser }}
      - SPRING_DATASOURCE_PASSWORD={{ $.DbPass }}
      {{- end }}

      {{- if .Features.HasMariaDB }}
      - SPRING_DATASOURCE_URL=jdbc:mariadb://mariadb:3306/{{ $.DbName }}?createDatabaseIfNotExist=true
      - SPRING_DATASOURCE_USERNAME={{ $.DbUser }}
      - SPRING_DATASOURCE_PASSWORD={{ $.DbPass }}
      {{- end }}

      {{- if .Features.HasMongo }}
      - SPRING_DATA_MONGODB_URI=mongodb://mongo:27017/{{ $.DbName }}
      {{- end }}

      {{- if .Features.HasCassandra }}
      - SPRING_DATA_CASSANDRA_CONTACT_POINTS=cassandra
      - SPRING_DATA_CASSANDRA_PORT=9042
      - SPRING_DATA_CASSANDRA_LOCAL_DATACENTER=datacenter1
      {{- end }}

      {{- if .Features.HasRedis }}
      - SPRING_DATA_REDIS_HOST=redis
      - SPRING_DATA_REDIS_PORT=6379
      {{- end }}

      {{- if .Features.HasElastic }}
      - SPRING_ELASTICSEARCH_URIS=http://elasticsearch:9200
      {{- end }}

      # --- Messaging ---
      {{- if .Features.HasRabbit }}
      - SPRING_RABBITMQ_HOST=rabbitmq
      - SPRING_RABBITMQ_PORT=5672
      - SPRING_RABBITMQ_USERNAME=guest
      - SPRING_RABBITMQ_PASSWORD=guest
      {{- end }}

      {{- if .Features.HasKafka }}
      # Важно: используем внутренний порт 29092
      - SPRING_KAFKA_BOOTSTRAP_SERVERS=kafka:29092
      {{- end }}

      {{- if .Features.HasActiveMQ }}
      - SPRING_ACTIVEMQ_BROKER_URL=tcp://activemq:61616
      {{- end }}

      # --- Discovery & Config ---
      {{- if .Features.HasEureka }}
      - EUREKA_CLIENT_SERVICEURL_DEFAULTZONE=http://eureka:8761/eureka/
      - EUREKA_INSTANCE_PREFER_IP_ADDRESS=true
      {{- end }}

      {{- if .Features.HasConsul }}
      - SPRING_CLOUD_CONSUL_HOST=consul
      - SPRING_CLOUD_CONSUL_PORT=8500
      {{- end }}

      {{- if .Features.HasConfigClient }}
      # Предполагаем, что Config Server называется 'config-server' (стандарт)
      - SPRING_CONFIG_IMPORT=optional:configserver:http://config-server:8888
      {{- end }}

      # --- Observability & Security ---
      {{- if .Features.HasZipkin }}
      - MANAGEMENT_ZIPKIN_TRACING_ENDPOINT=http://zipkin:9411/api/v2/spans
      {{- end }}

      {{- if .Features.HasVault }}
      - SPRING_CLOUD_VAULT_URI=http://vault:8200
      - SPRING_CLOUD_VAULT_TOKEN=root # Dev token
      {{- end }}

    depends_on:
      {{- if .Features.HasPostgres }}
      postgres:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasMySQL }}
      mysql:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasMariaDB }}
      mariadb:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasMongo }}
      mongo:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasRedis }}
      redis:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasElastic }}
      elasticsearch:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasKafka }}
      kafka:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasRabbit }}
      rabbitmq:
        condition: service_healthy
      {{- end }}
      {{- if .Features.HasEureka }}
      eureka:
        condition: service_started
      {{- end }}
      {{- if .Features.HasConsul }}
      consul:
        condition: service_started
      {{- end }}
      {{- if .Features.HasConfigClient }}
      # Если есть конфиг сервер, он должен стартовать первым
      # config-server: 
      #   condition: service_started
      {{- end }}
{{- end }}

  # ==========================================
  # Infrastructure Services
  # ==========================================

{{- if $.Features.HasPostgres }}
  postgres:
    image: postgres:16-alpine
    container_name: ghost_postgres
    networks:
      - ghost_network
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: {{ $.DbName }}
      POSTGRES_USER: {{ $.DbUser }}
      POSTGRES_PASSWORD: {{ $.DbPass }}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {{ $.DbUser }} -d {{ $.DbName }}"]
      interval: 5s
      timeout: 5s
      retries: 5
{{- end }}

{{- if $.Features.HasMySQL }}
  mysql:
    image: mysql:8.0.39
    container_name: ghost_mysql
    networks:
      - ghost_network
    ports:
      - "3306:3306"
    environment:
      MYSQL_DATABASE: {{ $.DbName }}
      MYSQL_USER: {{ $.DbUser }}
      MYSQL_PASSWORD: {{ $.DbPass }}
      MYSQL_ROOT_PASSWORD: root
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin" ,"ping", "-h", "localhost"]
      interval: 5s
      timeout: 5s
      retries: 5
{{- end }}

{{- if $.Features.HasMariaDB }}
  mariadb:
    image: mariadb:10.11.10
    container_name: ghost_mariadb
    networks:
      - ghost_network
    ports:
      - "3306:3306"
    environment:
      MARIADB_DATABASE: {{ $.DbName }}
      MARIADB_USER: {{ $.DbUser }}
      MARIADB_PASSWORD: {{ $.DbPass }}
      MARIADB_ROOT_PASSWORD: root
    volumes:
      - mariadb_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "healthcheck.sh", "--connect", "--innodb_initialized"]
      interval: 5s
      timeout: 5s
      retries: 5
{{- end }}

{{- if $.Features.HasMongo }}
  mongo:
    image: mongo:6.0.20
    container_name: ghost_mongo
    networks:
      - ghost_network
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 10s
      retries: 5
{{- end }}

{{- if $.Features.HasCassandra }}
  cassandra:
    image: cassandra:4.1.7
    container_name: ghost_cassandra
    networks:
      - ghost_network
    ports:
      - "9042:9042"
    environment:
      - CASSANDRA_CLUSTER_NAME=GhostCluster
      - CASSANDRA_DC=datacenter1
      - CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch
    volumes:
      - cassandra_data:/var/lib/cassandra
    healthcheck:
      test: ["CMD-SHELL", "cqlsh -e 'DESCRIBE CLUSTER'"]
      interval: 20s
      timeout: 10s
      retries: 5
{{- end }}

{{- if $.Features.HasRedis }}
  redis:
    image: redis:7.4.2-alpine
    container_name: ghost_redis
    networks:
      - ghost_network
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
{{- end }}

{{- if $.Features.HasElastic }}
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.17.23
    container_name: ghost_elastic
    networks:
      - ghost_network
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m" # Ограничиваем память, чтобы Docker не умер
    ports:
      - "9200:9200"
    volumes:
      - elastic_data:/usr/share/elasticsearch/data
    healthcheck:
      test: ["CMD-SHELL", "curl -s http://localhost:9200/_cluster/health | grep -q 'status.*green\\|status.*yellow'"]
      interval: 10s
      timeout: 10s
      retries: 5
{{- end }}

{{- if $.Features.HasRabbit }}
  rabbitmq:
    image: rabbitmq:3.13.16-management
    container_name: ghost_rabbitmq
    networks:
      - ghost_network
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # GUI
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 10s
      timeout: 10s
      retries: 5
{{- end }}

{{- if $.Features.HasActiveMQ }}
  activemq:
    image: apache/activemq-classic:5.18.4
    container_name: ghost_activemq
    networks:
      - ghost_network
    ports:
      - "61616:61616" # TCP
      - "8161:8161"   # Web Console
{{- end }}

{{- if $.Features.HasKafka }}
  zookeeper:
    image: confluentinc/cp-zookeeper:7.7.3
    container_name: ghost_zookeeper
    networks:
      - ghost_network
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.7.3
    container_name: ghost_kafka
    networks:
      - ghost_network
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      # Слушаем на 29092 для внутренних контейнеров и на 9092 для хоста (разработчика)
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    healthcheck:
      test: nc -z localhost 9092 || exit 1
      interval: 10s
      timeout: 5s
      retries: 5
{{- end }}

{{- if $.Features.HasEureka }}
  eureka:
    image: springcloud/eureka:3.1.2
    container_name: ghost_eureka
    networks:
      - ghost_network
    ports:
      - "8761:8761"
{{- end }}

{{- if $.Features.HasConsul }}
  consul:
    image: hashicorp/consul:1.19.3
    container_name: ghost_consul
    networks:
      - ghost_network
    ports:
      - "8500:8500"
{{- end }}

{{- if $.Features.HasZipkin }}
  zipkin:
    image: openzipkin/zipkin
    container_name: ghost_zipkin
    networks:
      - ghost_network
    ports:
      - "9411:9411"
{{- end }}

{{- if $.Features.HasVault }}
  vault:
    image: hashicorp/vault:latest
    container_name: ghost_vault
    networks:
      - ghost_network
    ports:
      - "8200:8200"
    environment:
      VAULT_DEV_ROOT_TOKEN_ID: root
    cap_add:
      - IPC_LOCK
{{- end }}

{{- if $.Features.HasPrometheus }}
  prometheus:
    image: prom/prometheus
    container_name: ghost_prometheus
    networks:
      - ghost_network
    ports:
      - "9090:9090"
    # Для полноценной работы нужно монтировать prometheus.yml,
    # но пока просто поднимаем контейнер
{{- end }}

volumes:
{{- if $.Features.HasPostgres }}
  postgres_data:
    driver: local
{{- end }}

{{- if $.Features.HasMySQL }}
  mysql_data:
    driver: local
{{- end }}

{{- if $.Features.HasMariaDB }}
  mariadb_data:
    driver: local
{{- end }}

{{- if $.Features.HasMongo }}
  mongo_data:
    driver: local
{{- end }}

{{- if $.Features.HasCassandra }}
  cassandra_data:
    driver: local
{{- end }}

{{- if $.Features.HasElastic }}
  elastic_data:
    driver: local
{{- end }}