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
FROM eclipse-temurin:%s-jdk-alpine AS builder
WORKDIR /app

# Копируем файлы maven wrapper и настройки
COPY .mvn/ .mvn
COPY mvnw pom.xml ./

# Скачиваем зависимости (чтобы закэшировать этот слой)
# Ускоряет повторные сборки, если pom.xml не менялся
RUN ./mvnw dependency:go-offline

# Копируем исходный код и собираем приложение
COPY src ./src
RUN ./mvnw clean package -DskipTests

# Разделяем fat-jar на слои (Layertools)
# Это фича Spring Boot, позволяющая отделить библиотеки от кода
RUN java -Djarmode=layertools -jar target/*.jar extract

# --- Этап 2: Запуск (Runtime) ---
# Используем JRE (только среда выполнения), а не JDK.
# Alpine - очень легкая версия Linux.
FROM eclipse-temurin:%s-jre-alpine

WORKDIR /app

# Создаем пользователя с ограниченными правами
RUN addgroup -S spring && adduser -S spring -G spring

# Копируем слои из этапа сборки
# Порядок важен! Сначала редко меняющиеся (зависимости), потом часто меняющиеся (код).
COPY --from=builder /app/dependencies/ ./
COPY --from=builder /app/spring-boot-loader/ ./
COPY --from=builder /app/snapshot-dependencies/ ./
COPY --from=builder /app/application/ ./

# Переключаемся на безопасного пользователя
USER spring

EXPOSE %d

# Запускаем приложение через JarLauncher (оптимизировано для слоев)
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
