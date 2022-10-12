# terraform-small-lab
Terraform Small Lab

# Install Terrform
https://www.terraform.io/downloads

```
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```
# INSTALL AWS CLI  
https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

```
curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
sudo installer -pkg AWSCLIV2.pkg -target /

# validate installation
which aws
aws --version

#configure default aws access (key / secret)
aws configure
```

# INSTALL PHYTON FLASK 
``` 
cd api

# maybe you need to use pip instead of pip3
pip3 install Flask
flask --app api run
```