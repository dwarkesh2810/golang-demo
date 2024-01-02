# Golang Demo

# Beego Demo

**app.conf File Need To Add As Shown below**
app.conf file is located at conf/app.conf

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
```

**App.env FIle Need To Add As Shown below**
Make sure change the variables value as you required

```
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

BASE_UPLOAD_PATH="./assets/uploads/"

CONN="user=[dbuser] password=[dbpassword] dbname=[dbname] host=localhost sslmode=disable"

JWT_SEC_KEY=your_unique
```

**Usefull CLI Commands**
Given below commands are use to run the project.
make sure you have install postgresql for database

```
go mod tidy
bee migrate -driver=[dbdriver] -conn="[dbdriver]://[dbuser]:[dbpassword]@localhost:5432/[dbname]?sslmode=disable"
bee run -downdoc=true -gendoc=true
```
