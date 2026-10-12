FROM node:26-alpine3.23 as builder

WORKDIR /app

COPY package*.json .

RUN npm install

COPY . .

EXPOSE 3000

RUN npm run build


FROM node:24-alpine as runner
WORKDIR /app /app
ENTRYPOINT [ "npm" , "start" ]
