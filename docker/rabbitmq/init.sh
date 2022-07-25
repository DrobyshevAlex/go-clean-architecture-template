#!/bin/sh

# Create Rabbitmq user
( rabbitmqctl wait --timeout 60 $DOCKER_RABBITMQ_PID_FILE ; \
rabbitmqctl add_user $DOCKER_RABBITMQ_USER $DOCKER_RABBITMQ_PASSWORD 2>/dev/null ; \
rabbitmqctl set_user_tags $DOCKER_RABBITMQ_USER administrator ; \
rabbitmqctl set_permissions -p / $DOCKER_RABBITMQ_USER  ".*" ".*" ".*" ; \
echo "*** User '$DOCKER_RABBITMQ_USER' with password '$DOCKER_RABBITMQ_PASSWORD' completed. ***" ; \
echo "*** Log in the WebUI at port 15672 (example: http:/localhost:15672) ***") &

# $@ is used to pass arguments to the rabbitmq-server command.
# For example if you use it like this: docker run -d rabbitmq arg1 arg2,
# it will be as you run in the container rabbitmq-server arg1 arg2
rabbitmq-server $@
