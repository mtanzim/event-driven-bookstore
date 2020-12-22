bin/zookeeper-server-start.sh config/zookeeper.properties &
bin/kafka-server-start.sh config/server.properties &
bin/kafka-topics.sh --create --topic cart-payment --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic cart-payment-processed --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic cart-shipment --bootstrap-server localhost:9092
