# install nats-server
echo "Installing nats-server..."
mkdir nats && cd nats
curl -L https://github.com/nats-io/nats-server/releases/download/v2.11.9/nats-server-v2.11.9-linux-amd64.tar.gz -o nats-server-v2.11.9-linux-amd64.tar.gz
tar -xzf nats-server-v2.11.9-linux-amd64.tar.gz
mv nats-server-v2.11.9-linux-amd64/nats-server /usr/local/bin/nats-server

# install nats-cli
curl -L https://github.com/nats-io/natscli/releases/download/v0.2.4/nats-0.2.4-linux-amd64.zip -o nats-0.2.4-linux-amd64.zip
unzip nats-0.2.4-linux-amd64.zip
mv nats-0.2.4-linux-amd64/nats /usr/local/bin/nats

# Install Redis stack
echo "Installing Redis stack..."
sudo apt-get install lsb-release curl gpg
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
sudo chmod 644 /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install redis-stack-server