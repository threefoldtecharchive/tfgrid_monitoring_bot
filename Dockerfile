FROM node:16-alpine AS Build
WORKDIR /app
COPY package.json ./
RUN yarn install

FROM node:16-alpine
COPY --from=Build /app/node_modules ./node_modules
COPY package.json app.ts .env ./
CMD ["yarn", "ts-node","app.ts"]
