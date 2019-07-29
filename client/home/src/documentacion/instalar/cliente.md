---
title: "Instalar Cliente"
date: "2019-27-02"
---

siguiendo las instrucciones de jonatan mircha YOUTUBE

https://www.youtube.com/watch?v=s_mNK_lg2Jw&list=PLvq-jIkSeTUY3gY-ptuqkNEXZHsNwlkND&index=67

## Instalar NPM
```bash
sudo apt-get update
sudo apt-get install npm
```

## Instalar GIT
```bash
sudo apt-get update
sudo apt-get install git
```

## Instalar NVM
Instalar nodejs con NVM

```bash
sudo apt-get install build-essential libssl-dev
```

ir a [repositorio](https://github.com/creationix/nvm) para la instalación de NVM

```bash
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.34.0/install.sh | bash
source ~/.bashrc
nvm ls-remote
nvm install 10.15.2
nvm use 10.15.2
nvm ls
```

## Instalar nginx
```bash
sudo apt-get update
sudo apt-get install nginx
```

## Configurar Proxy
Configurar un Servidor Proxy
```bash
sudo vi /etc/nginx/sites-available/default
```

```nginx
server {
    listen 80;
    server_name example.com;
    location / {
        proxy_pass http://APP_PRIVATE_IP_ADDRESS:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
sudo service nginx restart
```
### Proxy Auxiliar
(no es recomendable) en caso de no tener un dominio se puede configurar la tarjeta de red en ubuntu SOLO EN DEV
```bash
sudo iptables -t nat -A PREROUTING -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3102
```
```bash
sudo service nginx restart
```

## Instalar PM2
Instalación i configuración PM2 para mantener ejecutando app nodejs 

```bash
npm install pm2 -g
```

```bash
pm2 start npm --name "{app_name}" -- run {script_name}
```

```bash
┌──────────────────────────────┬────┬─────────┬──────┬───────┬────────┬─────────┬────────┬──────┬───────────┬──────┬──────────┐
│ App name                     │ id │ version │ mode │ pid   │ status │ restart │ uptime │ cpu  │ mem       │ user │ watching │
├──────────────────────────────┼────┼─────────┼──────┼───────┼────────┼─────────┼────────┼──────┼───────────┼──────┼──────────┤
│ institute_client_admission   │ 4  │ 0.34.0  │ fork │ 15076 │ online │ 0       │ 3h     │ 0.1% │ 45.6 MB   │ yoel │ disabled │
│ institute_client_certificate │ 5  │ 0.34.0  │ fork │ 15194 │ online │ 0       │ 3h     │ 0.1% │ 45.7 MB   │ yoel │ disabled │
│ institute_client_librarie    │ 6  │ 0.34.0  │ fork │ 15321 │ online │ 0       │ 3h     │ 0.2% │ 45.6 MB   │ yoel │ disabled │
│ institute_client_messenger   │ 10 │ 0.34.0  │ fork │ 17085 │ online │ 0       │ 3h     │ 0.1% │ 47.1 MB   │ yoel │ disabled │
│ institute_client_monitoring  │ 7  │ 0.34.0  │ fork │ 15425 │ online │ 0       │ 3h     │ 0.1% │ 45.7 MB   │ yoel │ disabled │
│ institute_client_review      │ 8  │ 0.34.0  │ fork │ 15510 │ online │ 0       │ 3h     │ 0.1% │ 45.0 MB   │ yoel │ disabled │
│ institute_client_student     │ 9  │ 0.34.0  │ fork │ 15627 │ online │ 0       │ 3h     │ 0.1% │ 46.1 MB   │ yoel │ disabled │
│ institute_doc                │ 0  │ 0.34.0  │ fork │ 13110 │ online │ 0       │ 4h     │ 0.1% │ 48.2 MB   │ yoel │ disabled │
└──────────────────────────────┴────┴─────────┴──────┴───────┴────────┴─────────┴────────┴──────┴───────────┴──────┴──────────┘
```
