/home/mokammel/Applications/kafka-2.6.0-src/bin/zookeeper-server-start.sh config/zookeeper.properties &
/home/mokammel/Applications/kafka-2.6.0-src/bin/kafka-server-start.sh config/server.properties &
/home/mokammel/Applications/kafka-2.6.0-src/bin/kafka-topics.sh --create --topic cart-payment --bootstrap-server localhost:9092
/home/mokammel/Applications/kafka-2.6.0-src/bin/kafka-topics.sh --create --topic cart-payment-processed --bootstrap-server localhost:9092
/home/mokammel/Applications/kafka-2.6.0-src/bin/kafka-topics.sh --create --topic cart-shipment --bootstrap-server localhost:9092
