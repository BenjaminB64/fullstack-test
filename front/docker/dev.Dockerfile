# Use an official Node.js, and it should be version 16 and above
FROM node:21-alpine3.18
# Set the working directory in the container
WORKDIR /app
# Install app dependencies using PNPM
RUN npm install -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --force
# Expose the app
EXPOSE 3000
# Start the application
CMD ["pnpm", "run", "dev"]