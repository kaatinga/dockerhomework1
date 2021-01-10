
1. Rolling update обновляет постепенно, под за подом. 
   
```
m.onishenko@painter manifests % kubectl apply -f deployment.yaml      
deployment.apps/homework3 configured
m.onishenko@painter manifests % kubectl get pods                
NAME                         READY   STATUS              RESTARTS   AGE
homework3-64b6c8dbb-78862    1/1     Terminating         0          54s
homework3-64b6c8dbb-s7fdg    1/1     Running             0          54s
homework3-6f97668985-gc7rn   0/1     ContainerCreating   0          4s
homework3-6f97668985-xjdhd   0/1     ContainerCreating   0          3s
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS              RESTARTS   AGE
homework3-64b6c8dbb-78862    0/1     Terminating         0          61s
homework3-64b6c8dbb-s7fdg    1/1     Running             0          61s
homework3-6f97668985-gc7rn   0/1     ContainerCreating   0          11s
homework3-6f97668985-xjdhd   0/1     Running             0          10s
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
homework3-64b6c8dbb-s7fdg    1/1     Running   0          65s
homework3-6f97668985-gc7rn   0/1     Running   0          15s
homework3-6f97668985-xjdhd   0/1     Running   0          14s
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS        RESTARTS   AGE
homework3-64b6c8dbb-s7fdg    1/1     Terminating   0          71s
homework3-6f97668985-gc7rn   1/1     Running       0          21s
homework3-6f97668985-xjdhd   1/1     Running       0          20s
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS        RESTARTS   AGE
homework3-64b6c8dbb-s7fdg    0/1     Terminating   0          76s
homework3-6f97668985-gc7rn   1/1     Running       0          26s
homework3-6f97668985-xjdhd   1/1     Running       0          25s
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
homework3-6f97668985-gc7rn   1/1     Running   0          81s
homework3-6f97668985-xjdhd   1/1     Running   0          80s
```

2. Fixed update я реализовал указав тип стратегии recreate и удалив все другие настроки стратегии так как они не работают, файл deployment_fixed.yaml.

результат:
```
m.onishenko@painter manifests % kubectl apply -f deployment_fixed.yaml
deployment.apps/homework3 configured
m.onishenko@painter manifests % kubectl get pods                      
NAME                         READY   STATUS        RESTARTS   AGE
homework3-6f97668985-gc7rn   1/1     Terminating   0          12m
homework3-6f97668985-xjdhd   1/1     Terminating   0          12m
m.onishenko@painter manifests % kubectl get pods
NAME                         READY   STATUS        RESTARTS   AGE
homework3-6f97668985-xjdhd   0/1     Terminating   0          13m
m.onishenko@painter manifests % kubectl get pods
NAME                        READY   STATUS              RESTARTS   AGE
homework3-64b6c8dbb-5262x   0/1     ContainerCreating   0          0s
homework3-64b6c8dbb-sgnml   0/1     ContainerCreating   0          0s
```
Контейнеры уничтожаются сразу и после того как уничтожен последний, запускаются новые.

3. Blue-Green Deployment. Я подправил deployment_fixed так чтобы наименование селектора поменялось. Теперь я имею два деплоймента, но только 1 из них работает так как сервис повернут на старый селектор homework3.

Редактирую сервис:
```
m.onishenko@painter manifests % kubectl patch service homework3-srv -p '{"spec":{"selector":{"app": "homework3v1"}}}'
service/homework3-srv patched
```
Проверяю что поменялось:
```
m.onishenko@painter manifests % kubectl describe service homework3-srv
Name:                     homework3-srv
Namespace:                c2-mo
Labels:                   <none>
Annotations:              <none>
Selector:                 app=homework3v1
Type:                     NodePort
IP Families:              <none>
IP:                       10.254.21.176
IPs:                      <none>
Port:                     http  9090/TCP
TargetPort:               8080/TCP
NodePort:                 http  30118/TCP
Endpoints:                10.100.226.143:8080,10.100.231.56:8080
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```
Теперь я обращаюсь к новой ручке http://homework3.host/version

Вижу что работает сегодняшняя сборка с версией от 10 числа:
Build Time: `2021 +01 +10-08:06:36`

Патчу сервис обратно:
```
m.onishenko@painter manifests % kubectl patch service homework3-srv -p '{"spec":{"selector":{"app": "homework3"}}}'  
service/homework3-srv patched
```
Обращаюсь к ручке http://homework3.host/version

Вижу `Build Time: 2021 +01 +07-09:53:46` от версии :latest

4. Canary Deployment. Я сделал еще 1 файл: deployment_canary.yaml в котором все теги одинаковые кроме имени деплоймента.
```
m.onishenko@painter manifests % kubectl apply -f deployment-canary.yaml                                              
deployment.apps/homework3-v1 created
m.onishenko@painter manifests % kubectl get pods                                                                   
NAME                           READY   STATUS              RESTARTS   AGE
homework3-6f97668985-tp2kx     1/1     Running             0          40m
homework3-6f97668985-tq2s7     1/1     Running             0          40m
homework3-v1-64b6c8dbb-j4p4f   0/1     ContainerCreating   0          8s
homework3-v1-64b6c8dbb-rz49z   0/1     ContainerCreating   0          8s
homework3v1-5f9d4975c6-nx968   1/1     Running             0          57m
homework3v1-5f9d4975c6-vp5s4   1/1     Running             0          57m
```
Самое главное селектор тот же самый. Теперь сервис будет выбирать под из двух деплойментов примерно в равной пропорции. Так и есть: http://homework3.host/version теперь при каждом обновлении выдаёт разную версию. Теперь я могу сокращать или увеличивать количество подов в каждом деплойменте чтобы плавно переходить от версии к версии:
```
m.onishenko@painter manifests % kubectl scale deployment/homework3-v1 --replicas=4
deployment.apps/homework3-v1 scaled
m.onishenko@painter manifests % kubectl get deployment homework3-v1               
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
homework3-v1   2/4     4            2           3m51s
m.onishenko@painter manifests % kubectl get deployment homework3-v1
NAME           READY   UP-TO-DATE   AVAILABLE   AGE
homework3-v1   4/4     4            4           3m56s
```
Теперь при каждом обновлении http://homework3.host/version чаще показывается именно версия от 10 числа.

Всё запущено и работает по адресу http://homework3.host если подправить /etc/hosts