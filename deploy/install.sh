echo "Preparing XCPC Team Registration API Service..."

echo "Initializing database..."

echo -n "Please input the host of database [127.0.0.1]: "
read host
host=${host:-"127.0.0.1"}

echo -n "Please input the port of database [5432]: "
read port
port=${port:-"5432"}

echo -n "Please input the username of database [postgres]: "
read usrname
usrname=${usrname:-"postgres"}

psql -h $host -p $port -U $usrname -W -f ./sql/init.sql

if [ "$?" -ne 0 ]
then
    echo "Database initialization failed"
    exit 0    
fi

echo "Database is initialized successfully"

echo "Initializing Redis..."

echo "flushall" | redis-cli

echo "Redis is initialized successfully"

echo "Building server..."

go build -o ../ ../

if [ "$?" -ne 0 ]
then
    echo "Build server failed"
    exit 0    
fi

echo "Build server finished"

echo "Initializing root user..."

arr=$(python3 ./genpwd.py)
pwdtoken=$(echo $arr | awk '{print $1}')
rootpwd=$(echo $arr | awk '{print $2}')
token=$(echo $arr | awk '{print $3}')

echo -n "Input the root's email (for login): "

read email

echo "$pwdtoken $email" | ../xcpc-team-reg -i -c ../conf/config.yaml

if [ "$?" -ne 0 ]
then
    echo "Initialize root user failed"
    exit 0    
fi

echo "Root user password is: $rootpwd"
echo "Website token is: $token"
echo "run '../xcpc-team-reg -c ../conf/config.yaml' to start server"
