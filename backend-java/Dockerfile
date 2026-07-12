# ---- Build ----
FROM maven:3-eclipse-temurin-25 AS build
WORKDIR /app
COPY pom.xml .
RUN mvn -B dependency:go-offline
COPY src ./src
RUN mvn -B clean package -DskipTests

# ---- Runtime ----
FROM eclipse-temurin:25-jre
WORKDIR /app
COPY --from=build /app/target/portal-empleado-1.0.0.jar app.jar

# La base de datos SQLite vive en el volumen montado en /data; las migraciones
# se montan en /migrations desde shared/migrations (ver docker-compose).
ENV PORTAL_DB_PATH=/data/portal.db
ENV MIGRATIONS_PATH=/migrations
ENV SERVER_PORT=8080
EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
