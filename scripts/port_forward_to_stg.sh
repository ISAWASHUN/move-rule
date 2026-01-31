#!/bin/sh
INSTANCE_ID=$(aws ec2 describe-instances \
	--filters "Name=tag:Name,Values=garbage-category-rule-quiz-bastion-stg" "Name=instance-state-name,Values=running" \
	--query "Reservations[0].Instances[0].InstanceId" \
	--output text \
	--profile garbage-category-rule-quiz-stg)

aws ssm start-session \
	--target $INSTANCE_ID \
	--document-name AWS-StartPortForwardingSessionToRemoteHost \
	--parameters '{"host":["garbage-category-rule-quiz-stg.cbks6uwysn53.ap-northeast-1.rds.amazonaws.com"],"portNumber":["3306"], "localPortNumber":["13306"]}' \
	--profile garbage-category-rule-quiz-stg
