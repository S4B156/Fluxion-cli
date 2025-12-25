package generator

import (
	"fmt"
	"project/pkg/colors"
	"os"
	"path/filepath"
	"project/models"
	"strconv"
	"strings"
)

func GenerateDockerfile(rootPath string, info models.Project, port int) {
	launch := ""
	if !Is31OrOlder(info.Parent.SpringVersion) {
		launch = ".launch"
	}
	
	dockerfileContent := fmt.Sprintf(
`# --- Этап 1: Сборка (Builder) ---
# Используем образ с предустановленным Maven. 
# Это надежнее, чем полагаться на mvnw в старых проектах.
FROM maven:3.9-eclipse-temurin-%s-alpine AS builder
WORKDIR /app

# Копируем только файл зависимостей
COPY pom.xml ./

# Скачиваем зависимости (кэширование Docker слоя)
RUN mvn dependency:go-offline -B

# Копируем исходный код
COPY src ./src

# Собираем приложение
RUN mvn clean package -DskipTests

# Разделяем fat-jar на слои (Layertools)
# (Тут всё ок, оставляем как было)
RUN java -Djarmode=layertools -jar target/*.jar extract

# --- Этап 2: Запуск (Runtime) ---
# Используем легковесную JRE для запуска
FROM eclipse-temurin:%s-jre-alpine

WORKDIR /app

# Создаем пользователя (Security best practice)
RUN addgroup -S spring && adduser -S spring -G spring

# Копируем слои из билдера
COPY --from=builder /app/dependencies/ ./
COPY --from=builder /app/spring-boot-loader/ ./
COPY --from=builder /app/snapshot-dependencies/ ./
COPY --from=builder /app/application/ ./

USER spring

EXPOSE %d

ENTRYPOINT ["java", "org.springframework.boot.loader%s.JarLauncher"]
`, info.Properties.JavaVersion, info.Properties.JavaVersion, port, launch)

	filePath := filepath.Join(rootPath, "Dockerfile")
	err := os.WriteFile(filePath, []byte(dockerfileContent), 0644)
	if err != nil {
		colors.Error.Println("Error writing Dockerfile:", err)
	} else {
		colors.Docker.Printf("[ DONE ] Dockerfile created at %s\n", filePath)
	}
}

func Is31OrOlder(v string) bool {
	parts := strings.Split(v, ".")
	if len(parts) == 0 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major > 3 {
		return false
	}
	if major < 3 {
		return true
	}
	minor, err := strconv.Atoi(parts[1])
	return err == nil && minor <= 1
}
