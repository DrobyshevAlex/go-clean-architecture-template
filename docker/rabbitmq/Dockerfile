FROM rabbitmq:3-management

# Define environment variables.
ENV DOCKER_RABBITMQ_USER dev
ENV DOCKER_RABBITMQ_USER dev
ENV DOCKER_RABBITMQ_USER /var/lib/rabbitmq/mnesia/rabbitmq

ADD init.sh /init.sh
RUN chmod +x /init.sh
EXPOSE 15672

# Define default command
CMD ["/init.sh"]
