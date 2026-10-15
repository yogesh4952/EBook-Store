FROM node:latest AS builder

WORKDIR /app 

COPY package*.json ./ 
RUN npm install

EXPOSE 3000
COPY . .
RUN npm run build


FROM node:latest AS runner
WORKDIR /app 
ENV NODE_ENV=production
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

RUN useradd app
USER app

CMD ["npm","start"]
