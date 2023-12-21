# Golang Demo

**app.conf File Need To Add As Shown below**

```
appname = golang-demo
httpport = 8080
httpaddr = localhost
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true

#session settings
sessionOn=true
SessionGCMaxLifetime=60
SessionCookieLifeTime=60

# Enable i18n support
Enablei18n = true

#admin
EnableAdmin = true
AdminAddr = 0.0.0.0
AdminPort = 8000

[golang-demo]
dbdriver = dbdriver
dbusername = dbuser
dbpassword = dbpassword
dbname = dbname
dbhost = dbhost
dbport = dbport

RATELIMITER = 10
BLOCKTIME = 60 

EMAIL="xyz@gmail.com"
PASSWORD = xyzzz
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

**Usefull CLI Commands**

```
go mod tidy
bee migrate -driver=[dbdriver] -conn="[dbdriver]://[dbuser]:[dbpassword]@localhost:5432/dbname?sslmode=disable"
bee run -downdoc=true -gendoc=true
```
