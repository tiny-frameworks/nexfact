# Wir nutzen die JRE (Java Runtime Environment), nicht das volle JDK, um Platz zu sparen.
# "jammy" ist das stabile Ubuntu 22.04 LTS.
FROM eclipse-temurin:21-jre-jammy

# 1. Update und Installation der wichtigsten Schriften für PDFBox
RUN apt-get update && apt-get install -y --no-install-recommends \
    fontconfig \
    fonts-dejavu-core \
    fonts-liberation \
    && rm -rf /var/lib/apt/lists/*

# 2. (Optional) Font-Cache aktualisieren, damit Java die neuen Schriften sofort sieht
RUN fc-cache -f -v

# 3. Arbeitsverzeichnis anlegen
WORKDIR /app

# 4. Dein kompiliertes Java-Programm kopieren
COPY nexfact/target/zugferd-1.0.0.jar /app/processor.jar

ENV NEXFACT_MODE=watch
ENV WATCH_DIR=/data/in

# 5. Der Entrypoint (Argumente werden dann vom Watcher/Go durchgereicht)
ENTRYPOINT ["java", "-jar", "/app/processor.jar"]

# 6. OCI Metadata Labels
LABEL org.opencontainers.image.title="nexfact-zugferd"
LABEL org.opencontainers.image.description="Factur-X / ZUGFeRD document processing service" 
LABEL org.opencontainers.image.source="https://codeberg.org/tiny-frameworks/nexfact"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.authors="Georg Hagn <tiny-frameworks>"
LABEL org.opencontainers.image.url="https://codeberg.org/tiny-frameworks/nexfact"
LABEL org.opencontainers.image.documentation="https://codeberg.org/tiny-frameworks/nexfact/src/branch/main/README.md"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.vendor="tiny-frameworks"

