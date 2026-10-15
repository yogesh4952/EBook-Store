FROM node:latest

WORKDIR /app 

COPY package*.json ./ 
RUN npm install

EXPOSE 3000
COPY . .
RUN npm run build

RUN useradd app
USER app

CMD ["npm","start"]
