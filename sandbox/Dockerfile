FROM node:10
LABEL artifactId="skrop-sandbox"
LABEL maintainer="Zalando SE"

# Create app directory
WORKDIR /home/skrop/app

# Install app dependencies
COPY package*.json ./
COPY . .

RUN npm install
RUN echo "npm install - SUCCESS"
RUN npm run build
RUN echo "npm run build - SUCCESS"


ENV NODE_ENV production
EXPOSE 8080
CMD ["npm", "start"]
