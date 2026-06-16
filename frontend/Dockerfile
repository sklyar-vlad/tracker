# -- BUILD STAGE --
FROM node:24-alpine3.22 AS builder

WORKDIR /app

COPY package.json package-lock.json* ./

RUN --mount=type=cache,target=/root/.npm npm ci

COPY . . 

RUN npm run build

# -- RUN STAGE --
FROM nginx:1.28.0-alpine3.21 AS runner

COPY nginx.conf /etc/nginx/nginx.conf

COPY --chown=nginx:nginx --from=builder /app/dist /usr/share/nginx/html

EXPOSE 5173

ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]