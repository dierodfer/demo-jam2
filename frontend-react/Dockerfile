FROM node:22-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .

EXPOSE 5173
# Servidor de desarrollo de Vite: suficiente para la demo 100% local y permite
# leer VITE_API_BASE del entorno en tiempo de arranque.
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]
