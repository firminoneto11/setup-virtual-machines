package scripts

const Python string = `
	sudo apt update
	sudo apt upgrade -y
	sudo apt install software-properties-common build-essential -y
	sudo add-apt-repository ppa:deadsnakes/ppa -y
	sudo apt update
	sudo apt install python3.11 python3.11-venv -y
`

const Docker string = `
	sudo apt update
	sudo apt upgrade -y
	sudo apt remove docker docker-engine docker.io containerd runc
	sudo apt install curl apt-transport-https ca-certificates gnupg2 lsb-release -y
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
	echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
		https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
	sudo apt update
	sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y
`

const Nginx string = `
	sudo apt update
	sudo apt upgrade -y
	sudo apt install curl gnupg2 ca-certificates lsb-release ubuntu-keyring
	curl https://nginx.org/keys/nginx_signing.key | gpg --dearmor | sudo tee /usr/share/keyrings/nginx-archive-keyring.gpg >/dev/null
	echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/ubuntu ` + "`lsb_release -cs` " + `
	nginx" | sudo tee /etc/apt/sources.list.d/nginx.list
	echo -e "Package: *\nPin: origin nginx.org\nPin: release o=nginx\nPin-Priority: 900\n" | sudo tee /etc/apt/preferences.d/99nginx
	sudo apt update
	sudo apt install nginx
`

const NodejsAndYarn string = `
	sudo apt update
	sudo apt upgrade -y
	curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
	sudo apt install nodejs -y
	sudo apt update
	sudo apt upgrade -y
	curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
	echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
	sudo apt update
	sudo apt install yarn -y
`
