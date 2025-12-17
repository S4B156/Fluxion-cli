version: '3.9'  # Используем свежую версию

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
    environment:
      - SERVER_PORT={{ .Port }}
      - SPRING_PROFILES_ACTIVE=docker
      
      # --- Database Connection ---
      {{- if $.Features.HasPostgres }}
      - SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/{{ $.DbName }}
      - SPRING_DATASOURCE_USERNAME={{ $.DbUser }}
      - SPRING_DATASOURCE_PASSWORD={{ $.DbPass }}
      - SPRING_JPA_HIBERNATE_DDL_AUTO=update
      {{- else if $.Features.HasMySQL }}
      - SPRING_DATASOURCE_URL=jdbc:mysql://mysql:3306/{{ $.DbName }}?useSSL=false&allowPublicKeyRetrieval=true
      - SPRING_DATASOURCE_USERNAME={{ $.DbUser }}
      - SPRING_DATASOURCE_PASSWORD={{ $.DbPass }}
      {{- end }}
      {{- if $.Features.HasMongo }}
      - SPRING_DATA_MONGODB_URI=mongodb://mongo:27017/{{ $.DbName }}
      {{- end }}

      # --- Cache & Broker ---
      {{- if $.Features.HasRedis }}
      - SPRING_DATA_REDIS_HOST=redis
      - SPRING_DATA_REDIS_PORT=6379
      {{- end }}
      {{- if $.Features.HasRabbit }}
      - SPRING_RABBITMQ_HOST=rabbitmq
      - SPRING_RABBITMQ_PORT=5672
      {{- end }}
      {{- if $.Features.HasKafka }}
      - SPRING_KAFKA_BOOTSTRAP_SERVERS=kafka:29092
      {{- end }}
      
      # --- Discovery ---
      {{- if $.Features.HasEureka }}
      - EUREKA_CLIENT_SERVICEURL_DEFAULTZONE=http://eureka:8761/eureka/
      {{- end }}

    depends_on:
      {{- if $.Features.HasPostgres }}
      postgres:
        condition: service_healthy
      {{- end }}
      {{- if $.Features.HasMySQL }}
      mysql:
        condition: service_healthy
      {{- end }}
      {{- if $.Features.HasMongo }}
      mongo:
        condition: service_healthy
      {{- end }}
      {{- if $.Features.HasKafka }}
      kafka:
        condition: service_started
      {{- end }}
      {{- if $.Features.HasRabbit }}
      rabbitmq:
        condition: service_healthy
      {{- end }}
      {{- if $.Features.HasEureka }}
      eureka:
        condition: service_started
      {{- end }}
{{- end }}

  # ==========================================
  # Infrastructure
  # ==========================================

{{- if .Features.HasPostgres }}
  postgres:
    image: postgres:15-alpine
    container_name: ghost_postgres
    networks:
      - ghost_network
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: {{ .DbName }}
      POSTGRES_USER: {{ .DbUser }}
      POSTGRES_PASSWORD: {{ .DbPass }}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {{ .DbUser }}"]
      interval: 5s
      timeout: 5s
      retries: 5
{{- end }}

{{- if .Features.HasMySQL }}
  mysql:
    image: mysql:8.0
    container_name: ghost_mysql
    networks:
      - ghost_network
    ports:
      - "3306:3306"
    environment:
      MYSQL_DATABASE: {{ .DbName }}
      MYSQL_USER: {{ .DbUser }}
      MYSQL_PASSWORD: {{ .DbPass }}
      MYSQL_ROOT_PASSWORD: root
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin" ,"ping", "-h", "localhost"]
      interval: 5s
      timeout: 5s
      retries: 5
{{- end }}

{{- if .Features.HasMongo }}
  mongo:
    image: mongo:6.0
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

{{- if .Features.HasRedis }}
  redis:
    image: redis:alpine
    container_name: ghost_redis
    networks:
      - ghost_network
    ports:
      - "6379:6379"
{{- end }}

{{- if .Features.HasRabbit }}
  rabbitmq:
    image: rabbitmq:3-management
    container_name: ghost_rabbitmq
    networks:
      - ghost_network
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 10s
      timeout: 10s
      retries: 5
{{- end }}

{{- if .Features.HasKafka }}
  zookeeper:
    image: confluentinc/cp-zookeeper:7.3.0
    container_name: ghost_zookeeper
    networks:
      - ghost_network
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000

  kafka:
    image: confluentinc/cp-kafka:7.3.0
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
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
{{- end }}

{{- if .Features.HasElastic }}
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.17.9
    container_name: ghost_elastic
    networks:
      - ghost_network
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - elastic_data:/usr/share/elasticsearch/data
{{- end }}

{{- if .Features.HasEureka }}
  eureka:
    image: springcloud/eureka
    container_name: ghost_eureka
    networks:
      - ghost_network
    ports:
      - "8761:8761"
{{- end }}

volumes:
{{- if .Features.HasPostgres }}  postgres_data:{{ end }}
{{- if .Features.HasMySQL }}  mysql_data:{{ end }}
{{- if .Features.HasMongo }}  mongo_data:{{ end }}
{{- if .Features.HasElastic }}  elastic_data:{{ end }}