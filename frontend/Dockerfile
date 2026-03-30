# Stage 1: Build
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

ARG API_BASE_URL

# Set environment variables for the build
ENV VITE_API_URL=${API_BASE_URL}

RUN npm run build

# Stage 2: Serve with Nginx as Non-Root
FROM nginx:stable-alpine

RUN addgroup -S nginxgroup && adduser -S nginxuser -G nginxgroup

# Adjust permissions so nginxuser can run the process
# nginx needs access to these directories to manage cache and PIDs
RUN touch /var/run/nginx.pid && \
    chown -R nginxuser:nginxgroup /var/run/nginx.pid && \
    chown -R nginxuser:nginxgroup /var/cache/nginx && \
    chown -R nginxuser:nginxgroup /var/log/nginx && \
    chown -R nginxuser:nginxgroup /etc/nginx/conf.d

# Switch to non-root user
USER nginxuser

# Copy static files (ensure they are owned by the new user)
COPY --from=builder --chown=nginxuser:nginxgroup /app/dist /usr/share/nginx/html

# copy custom nginx config as a template; envsubst will substitute
# ${API_BASE_URL} at container start, leaving nginx variables ($uri etc.) intact
COPY --chown=nginxuser:nginxgroup nginx.conf /etc/nginx/conf.d/default.conf.template

EXPOSE 8001

CMD ["/bin/sh", "-c", "envsubst '${API_BASE_URL}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]