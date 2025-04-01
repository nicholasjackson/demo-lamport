FROM node:22.14.0

COPY . /app

WORKDIR /app/ui
RUN npm install

ENTRYPOINT ["npm", "run", "dev"]