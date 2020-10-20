FROM bitnami/redis:latest
COPY database/dbquizcontent/redis.conf /opt/bitnami/redis/mounted-etc/redis.conf
CMD ["redis-server","/opt/bitnami/redis/mounted-etc/redis.conf", "--appendonly", "yes"]
USER 1001
EXPOSE 6381