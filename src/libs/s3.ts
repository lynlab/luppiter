import S3 from "aws-sdk/clients/s3";

export default class S3Client {
  private s3: S3;
  private bucketName: string;

  constructor(bucketName: string) {
    this.s3 = new S3({ region: "ap-northeast-2" });
    this.bucketName = bucketName;
  }

  public read(key: string): Promise<S3.GetObjectOutput> {
    const params = {
      Bucket: this.bucketName,
      Key: key,
    };

    return new Promise((resolve, reject) => {
      this.s3.getObject(params, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  public list(prefix?: string, cursor?: string): Promise<S3.ListObjectsV2Output> {
    const params = {
      Bucket: this.bucketName,
      Delimiter: "/",
      Prefix: prefix,
      MaxKeys: 20,
      StartAfter: cursor,
    };

    return new Promise((resolve, reject) => {
      this.s3.listObjectsV2(params, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }

  public write(key: string, buffer: Buffer): Promise<S3.PutObjectOutput> {
    const params = {
      Bucket: this.bucketName,
      Key: key,
      Body: buffer,
    };

    return new Promise((resolve, reject) => {
      this.s3.putObject(params, (err, data) => {
        err ? reject(err) : resolve(data);
      });
    });
  }
}
