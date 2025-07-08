FROM node:24.2.0-slim

# Install dependencies and create a user with user group
USER root

RUN addgroup --gid 1042 vapusdata

RUN adduser -D -u 503 -G vapusdata vapusdata-webapp-app

RUN mkdir -p /src

WORKDIR /src

# Copy package files first for better caching
COPY package*.json ./

RUN npm install style-loader css-loader --save-dev

# Install dependencies
RUN npm ci --only=production

# Copy source code
COPY . .

RUN chown -R vapusdata-webapp-app:vapusdata /src

USER vapusdata-webapp-app

RUN npm run build

EXPOSE 9017

RUN date -u

CMD ["npm", "start", "--", "-p", "9017"]