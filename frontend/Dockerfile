FROM node:20-alpine

WORKDIR /app

# Копирование package.json
COPY package*.json ./

# Установка зависимостей
RUN npm ci

# Копирование исходников
COPY . .

# Сборка приложения
RUN npm run build

EXPOSE 3000

CMD ["npm", "start"]
