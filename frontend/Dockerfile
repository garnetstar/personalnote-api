FROM node:20-alpine AS build

WORKDIR /app

COPY app/package*.json ./
RUN npm install

COPY app ./

RUN npm run build

FROM node:20-alpine

WORKDIR /usr/share/app

RUN npm install -g serve

COPY --from=build /app/dist ./dist

EXPOSE 3000

CMD ["serve", "-s", "dist", "-l", "3000"]
