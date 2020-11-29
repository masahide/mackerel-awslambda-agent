import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as s3 from '@aws-cdk/aws-s3';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as sqs from '@aws-cdk/aws-sqs';
import { SqsEventSource } from '@aws-cdk/aws-lambda-event-sources';

//import ec2 = require("@aws-cdk/aws-ec2");

export type Environment = {
    DRYRUN:     string;
    CLI:        string;
    S3BUCKET:   string;
    QUEUENAME:  string;
    STATETABLE: string;
    CHECKSFUNC: string;
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
                name: 'ExecuteID',
                type: dynamodb.AttributeType.STRING
            },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
            timeToLiveAttribute: "ExpiredAt",
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

        // LambdaFunction
        props.environment.STATETABLE = stateTable.tableName;
        props.environment.QUEUENAME = queue.queueName; 
        const checksFunc = new lambda.Function(this, `checks`, {
            memorySize: 1024,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/checks/'),
            handler: 'lambda-function',
            timeout: cdk.Duration.seconds(120),
           // role: appRole,
            environment: props.environment,
        });
        cdk.Tags.of(checksFunc).add('Name', `${props.stackName}-checks-function`);

        props.environment.CHECKSFUNC = checksFunc.functionName;
        const invokerFunc = new lambda.Function(this, `invoker`, {
            memorySize: 512,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/invoker/'),
            handler: 'lambda-function',
            timeout: cdk.Duration.seconds(120),
           // role: appRole,
            environment: props.environment,
        });
        cdk.Tags.of(invokerFunc).add('Name', `${props.stackName}-invoker-function`);
        const senderFunc = new lambda.Function(this, `sender`, {
            memorySize: 512,
            runtime: lambda.Runtime.GO_1_X,
            code: new lambda.AssetCode('../.dist/sender/'),
            handler: 'lambda-function',
            timeout: cdk.Duration.seconds(30),
           // role: appRole,
            environment: props.environment,
        });
        cdk.Tags.of(senderFunc).add('Name', `${props.stackName}-sender-function`);
        senderFunc.addEventSource(new SqsEventSource(queue));

        // IAM Role
        stateTable.grantFullAccess(checksFunc);
        bucket.grantRead(checksFunc);
        queue.grantSendMessages(checksFunc);
        queue.grantConsumeMessages(senderFunc);
    }
}
