FROM node:20-alpine AS build

# make the 'app' folder the current working directory
WORKDIR /app

# install simple http server for serving static content
RUN npm install -g http-server

# copy both 'package.json' and 'package-lock.json' (if available)
COPY ./web/package*.json ./

# install project dependencies
RUN npm install

# copy project files and folders to the current working directory (i.e. 'app' folder)
COPY ./web .

# build app for production with minification
RUN npm run build

EXPOSE 8004
CMD [ "http-server", "dist", "-p", "8004"]