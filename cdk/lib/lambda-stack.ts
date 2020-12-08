import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as s3 from '@aws-cdk/aws-s3';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as sqs from '@aws-cdk/aws-sqs';
import * as iam from '@aws-cdk/aws-iam';
import * as s3deploy from '@aws-cdk/aws-s3-deployment'
import * as path  from 'path';
import { SqsEventSource  } from '@aws-cdk/aws-lambda-event-sources';
import { Rule, Schedule } from '@aws-cdk/aws-events';
import { LambdaFunction } from '@aws-cdk/aws-events-targets';


//import ec2 = require("@aws-cdk/aws-ec2");

export type Environment = {
    DRYRUN:      string;
    CLI:         string;
    HOSTNAME:    string;
    S3KEY:       string;
	ORGANIZATION:string; 

    QUEUEURL:   string;
    S3BUCKET:   string;
    STATETABLE: string;
    CHECKERFUNC: string;
}

export interface Props extends cdk.StackProps {
    stackName: string;
    appName:   string;
    envName:   string;
    environment:  Environment;
}


export class LambdaStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props: Props) {
        super(scope, id, props);

        // DynamoDB table
        const stateTable = new dynamodb.Table(this, `stateTable`, {
            partitionKey: {
                name: 'id',
                type: dynamodb.AttributeType.STRING
            },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            timeToLiveAttribute: "ttl",
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        });
        cdk.Tags.of(stateTable).add('Name', `${props.stackName}-stateTable`);

        // SQS
        const queue = new  sqs.Queue(this,`queue`,{
           visibilityTimeout: cdk.Duration.seconds(300),
           receiveMessageWaitTime: cdk.Duration.seconds(20),
        });

        // S3 bucket
        const  bucket = new s3.Bucket(this, `configBucket`,{
            removalPolicy: cdk.RemovalPolicy.DESTROY,
        });
        props.environment.S3BUCKET = bucket.bucketName
        cdk.Tags.of(bucket).add('Name', `${props.stackName}-bucket`);

        // deploy config file
        new s3deploy.BucketDeployment(this, `${props.stackName}-config`, {
            sources: [s3deploy.Source.asset('../config',{ exclude: ['*.example']})],
            destinationKeyPrefix: path.dirname(props.environment.S3KEY),
            destinationBucket: bucket,
            prune: false,
        });
        // LambdaFunction
        props.environment.STATETABLE = stateTable.tableName;
        props.environment.QUEUEURL = queue.queueUrl; 
        const checkerFunc = new lambda.Function(this, `checker`, {
            memorySize: 1024,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/checker/'),
            handler: 'checker',
            timeout: cdk.Duration.seconds(120),
            environment: props.environment,
            //tracing: lambda.Tracing.ACTIVE,

        });
        cdk.Tags.of(checkerFunc).add('Name', `${props.stackName}-checker-function`);

        props.environment.CHECKERFUNC = checkerFunc.functionName;
        const invokerFunc = new lambda.Function(this, `invoker`, {
            memorySize: 128,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/invoker/'),
            handler: 'invoker',
            timeout: cdk.Duration.seconds(120),
            environment: props.environment,
            //tracing: lambda.Tracing.ACTIVE,
        });
        cdk.Tags.of(invokerFunc).add('Name', `${props.stackName}-invoker-function`);
        const senderFunc = new lambda.Function(this, `sender`, {
            memorySize: 128,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/sender/'),
            handler: 'sender',
            timeout: cdk.Duration.seconds(30),
            environment: props.environment,
            //tracing: lambda.Tracing.ACTIVE,
        });
        cdk.Tags.of(senderFunc).add('Name', `${props.stackName}-sender-function`);
        senderFunc.addEventSource(new SqsEventSource(queue));

        // IAM Role
        stateTable.grantFullAccess(invokerFunc);
        bucket.grantRead(invokerFunc);
        checkerFunc.grantInvoke(invokerFunc);
        stateTable.grantFullAccess(checkerFunc);
        bucket.grantRead(checkerFunc);
        queue.grantSendMessages(checkerFunc);
        checkerFunc.addToRolePolicy(
            new iam.PolicyStatement({
                actions: [
                    "logs:StartQuery",
                    "logs:StopQuery",
                    "logs:GetQueryResults",
                ],
                resources: [ "*" ],
            }),
        )
        queue.grantConsumeMessages(senderFunc);
        bucket.grantRead(senderFunc);

        new Rule(this, 'ScheduleRule', {
            schedule: Schedule.rate(cdk.Duration.minutes(1)),
            targets: [new LambdaFunction(invokerFunc)],
        });
    }
}
