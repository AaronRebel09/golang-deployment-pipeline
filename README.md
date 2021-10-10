# Deployment Pipeline for Go Applications on AWS

This repository provides an easy-to-deploy pipeline for the development, testing, building and deployment of applications written in Go. Although this example is tailored to Go, it can be easily modified to deploy applications written in other languages too.

Services Used:
 
 * [AWS CodePipeline](https://aws.amazon.com/codepipeline/) for pipeline creation.
 * [AWS CodeBuild](https://aws.amazon.com/codebuild/) for testing and building your Go application(s).
 * [AWS CloudFormation](https://aws.amazon.com/cloudformation/) for deploying infrastructure (Infrastructure-as-Code).
 * [AWS CodeDeploy](https://aws.amazon.com/codedeploy/) for zero downtime deployments of your application(s). 

#### Install the CodeDeploy agent for Amazon Linux or RHEL
```
#!/bin/bash
sudo yum update
sudo yum install ruby
sudo yum install wget
cd /home/ec2-user
wget https://aws-codedeploy-us-east-1.s3.us-east-1.amazonaws.com/latest/install
chmod +x ./install
sudo ./install auto
sudo yum install -y python-pip
sudo pip install awscli
sudo amazon-linux-extras install java-openjdk11
```

#### CodePipeline
```
1. Fetch data from github
2. CodeBuild with linux machine
3. CodeDeploy to EC2Instance
```