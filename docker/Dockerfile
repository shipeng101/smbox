FROM golang:1.21-alpine AS builder
# ... 构建后端

FROM node:18 AS frontend-builder
# ... 构建前端

FROM alpine:latest
# 复制构建产物
COPY --from=builder /app/smbox /usr/local/bin/
COPY --from=frontend-builder /app/dist /usr/share/smbox/web
# ... 其他配置 