docker_build('example-go-image', '.', 
    dockerfile='Dockerfile')
k8s_yaml('deployments/k8s.yaml')
k8s_yaml('deployments/postgres.yaml')
k8s_resource('example-go', port_forwards=8000, resource_deps=['postgres1'])
k8s_resource('postgres1', labels=['subscriptions-postgres'], port_forwards=5432)
