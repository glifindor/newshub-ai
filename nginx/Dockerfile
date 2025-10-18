# Production Nginx image for NewsHub AI
# Copies optimized configuration and prepares runtime directories
FROM nginx:1.27-alpine

# Ensure required directories exist for logs and certbot challenges
RUN mkdir -p /var/log/nginx \
    && mkdir -p /var/www/certbot \
    && chown -R nginx:nginx /var/log/nginx /var/www/certbot

# Remove default site configuration shipped with the base image
RUN rm -f /etc/nginx/conf.d/default.conf

# Copy configuration files
COPY nginx.conf /etc/nginx/nginx.conf
COPY nginx.prod.conf /etc/nginx/nginx.prod.conf

# Expose standard HTTP and HTTPS ports
EXPOSE 80 443

# Run Nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]
