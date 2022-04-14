docker_build('local-go-image', '.',
    dockerfile='Dockerfile',
    )

docker_build('local-postgres-image', '.',
    dockerfile='Dockerfile-postgres')
k8s_yaml('./deployments/local/go.yaml')
k8s_yaml('./deployments/local/postgres.yaml')
k8s_yaml('./deployments/local/krakend/krakend-config.yml')
k8s_yaml('./deployments/local/krakend/krakend-deployment.yml')
k8s_yaml('./deployments/local/krakend/krakend-service.yml')
k8s_resource('go-app', labels=['subscriptions'], port_forwards=8000, resource_deps=['postgres'])
k8s_resource('krakend-deployment', labels=['api-gateway'], port_forwards=8010)
k8s_resource('postgres', labels=['subscriptions'], port_forwards=5432)

k8s_yaml('./deployments/local/kafdrop-deployment.yaml')
k8s_yaml('./deployments/local/kafdrop-service.yaml')
k8s_yaml('./deployments/local/zookeeper.yaml')
k8s_yaml('./deployments/local/zookeeper-service.yaml')
k8s_yaml('./deployments/local/kafka.yaml')
k8s_yaml('./deployments/local/kafka-service.yaml')
// k8s_resource('zookeeper', labels=['event-streaming'], port_forwards=2181)
// k8s_resource('kafka', labels=['event-streaming'], resource_deps=['zookeeper'], port_forwards=9092)
// k8s_resource('kafdrop', labels=['event-streaming'], resource_deps=['kafka', 'zookeeper'], port_forwards=9000)